package main

import (
	"agent/tools"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
)

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []tools.ToolDefinition
	logger         *ChatLogger // set per chat session in Run
}

// And here:
func NewAgent(
	client *anthropic.Client,
	getUserMessage func() (string, bool),
	tools []tools.ToolDefinition,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

func main() {
	client := anthropic.NewClient()

	var getUserMessage func() (string, bool)
	if len(os.Args) > 1 {
		// One-shot mode: ./agent "do the task"
		prompt := strings.Join(os.Args[1:], " ")
		sent := false
		getUserMessage = func() (string, bool) {
			if sent {
				return "", false
			}
			sent = true
			return prompt, true
		}
	} else {
		// Interactive chat mode, same as before
		scanner := bufio.NewScanner(os.Stdin)
		getUserMessage = func() (string, bool) {
			if !scanner.Scan() {
				return "", false
			}
			return scanner.Text(), true
		}
	}

	// Remember what the user asked, so the reviewer sees the original problem and not just the agent's retelling.
	userMessages := []string{}
	readUserMessage := getUserMessage
	getUserMessage = func() (string, bool) {
		message, ok := readUserMessage()
		if ok {
			userMessages = append(userMessages, message)
		}
		return message, ok
	}
	getProblem := func() string { return strings.Join(userMessages, "\n\n") }

	tools := []tools.ToolDefinition{
		tools.ReadFileDefinition,
		tools.ListFilesDefinition,
		tools.EditFileDefinition,
		tools.RunProgramDefinition,
		tools.CreateFileDefinition,
		tools.NewReviewWorkDefinition(&client, getProblem),
	}
	agent := NewAgent(&client, getUserMessage, tools)

	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func (a *Agent) Run(ctx context.Context) error {
	logger, err := NewChatLogger("chat_logs")
	if err != nil {
		return err
	}
	defer logger.Close()
	a.logger = logger
	fmt.Printf("Chat log: %s\n", logger.file.Name())

	conversation := []anthropic.MessageParam{}

	fmt.Println("Chat with Clyppi, the most intelligent super intelligence of them all (use 'ctrl-c' to quit)")

	readUserInput := true
	for {
		if readUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			a.logger.Log("user", "%s", userInput)
			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			a.logger.Log("error", "inference failed: %s", err.Error())
			return err
		}
		conversation = append(conversation, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mClyppi\u001b[0m: %s\n", content.Text)
				a.logger.Log("assistant", "%s", content.Text)
			case "tool_use":
				result := a.executeTool(content.ID, content.Name, content.Input)
				toolResults = append(toolResults, result)
			}
		}
		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}
		readUserInput = false
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
	}

	return nil
}

// ChatLogger writes one chat session to its own log file.
type ChatLogger struct {
	mu   sync.Mutex
	file *os.File
}

// NewChatLogger creates a new, uniquely named log file in dir for a single chat.
// Every call produces a distinct file, so each chat session gets its own log.
func NewChatLogger(dir string) (*ChatLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	name := fmt.Sprintf("chat-%s-%d.log", time.Now().Format("20060102-150405.000000000"), os.Getpid())
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat log: %w", err)
	}

	l := &ChatLogger{file: f}
	l.Log("system", "chat session started")
	return l, nil
}

// Log appends a timestamped, role-tagged line to the chat's log file.
func (l *ChatLogger) Log(role, format string, args ...any) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "[%s] %s: %s\n", time.Now().Format(time.RFC3339Nano), role, fmt.Sprintf(format, args...))
}

// Close closes the chat's log file.
func (l *ChatLogger) Close() error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}

func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var toolDef tools.ToolDefinition
	var found bool
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}

	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
	a.logger.Log("tool_call", "%s(%s)", name, input)
	response, err := toolDef.Function(input)
	if err != nil {
		fmt.Printf("\u001b[91mconsole\u001b[0m: %s\n", err.Error())
		a.logger.Log("tool_error", "%s: %s", name, err.Error())
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	fmt.Printf("\u001b[96mconsole\u001b[0m: %s\n", response)
	a.logger.Log("tool_result", "%s: %s", name, response)
	return anthropic.NewToolResultBlock(id, response, false)
}

const systemPrompt = `You are a coding agent that solves programming tasks by reading, creating, editing and running files.
When you feel finished with coding, always get your work reviewed with the review_work tool before you stop.
If the reviewer answers "CHANGES NEEDED", fix the problems and ask for a new review. Stop once the reviewer answers "APPROVED", or after at most 3 reviews.`

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.BetaManagedAgentsModelClaudeSonnet5_5,
		MaxTokens: int64(16384),
		System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
		Messages:  conversation,
		Tools:     anthropicTools,
	})
	return message, err
}
