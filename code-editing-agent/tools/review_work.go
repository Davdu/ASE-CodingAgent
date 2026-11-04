package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// NewReviewWorkDefinition returns a tool that lets the agent ask a second LLM to review its work.
// It needs the client, so unlike the other tools it is built with a constructor.
// getProblem returns the original problem from the user, which is passed to the reviewer as is.
func NewReviewWorkDefinition(client *anthropic.Client, getProblem func() string) ToolDefinition {
	return ToolDefinition{
		Name:        "review_work",
		Description: "Ask an independent reviewer LLM to check the work you have done. Use this before you consider a task finished. The reviewer is given the original problem automatically. Describe the task as you understood it, give any extra context the reviewer needs, summarize what you did, list the files you changed, and say what you want checked. The reviewer can read and list files itself to verify your work, but cannot change anything.",
		InputSchema: ReviewWorkInputSchema,
		Function: func(input json.RawMessage) (string, error) {
			return ReviewWork(client, getProblem(), input)
		},
	}
}

type ReviewWorkInput struct {
	Task     string   `json:"task" jsonschema_description:"The original task you were asked to solve"`
	Context  string   `json:"context,omitempty" jsonschema_description:"Optional extra context for the reviewer, e.g. required function signatures, constraints or decisions you made"`
	Summary  string   `json:"summary" jsonschema_description:"A summary of what you have done to solve the task"`
	Files    []string `json:"files" jsonschema_description:"Relative paths of the files you created or changed, as a starting point for the reviewer"`
	CheckFor string   `json:"check_for" jsonschema_description:"What you want the reviewer to check, e.g. correctness, edge cases or whether the task is fully solved"`
}

var ReviewWorkInputSchema = GenerateSchema[ReviewWorkInput]()

// The reviewer may only read, never change, the workspace.
var reviewerTools = []ToolDefinition{ReadFileDefinition, ListFilesDefinition}

const maxReviewerTurns = 20

const reviewerSystemPrompt = `You are a strict code reviewer checking the work of another coding agent.
You get the original problem from the user, the task as the agent understood it, optional extra context, the agent's summary of what it did, and a list of files it touched.
The original problem is the source of truth. Do not trust the summary. Use your read_file and list_files tools to inspect the files yourself and verify the work.
When you are done, answer with a verdict on the first line, either "APPROVED" or "CHANGES NEEDED", followed by a short list of concrete problems and how to fix them.`

func ReviewWork(client *anthropic.Client, problem string, input json.RawMessage) (string, error) {
	reviewWorkInput := ReviewWorkInput{}
	err := json.Unmarshal(input, &reviewWorkInput)
	if err != nil {
		return "", err
	}
	Usage.RecordReview()

	var prompt strings.Builder
	fmt.Fprintf(&prompt, "<original_problem>\n%s\n</original_problem>\n\n", problem)
	fmt.Fprintf(&prompt, "<task>\n%s\n</task>\n\n", reviewWorkInput.Task)
	fmt.Fprintf(&prompt, "<context>\n%s\n</context>\n\n", reviewWorkInput.Context)
	fmt.Fprintf(&prompt, "<summary>\n%s\n</summary>\n\n", reviewWorkInput.Summary)
	fmt.Fprintf(&prompt, "<check_for>\n%s\n</check_for>\n\n", reviewWorkInput.CheckFor)
	fmt.Fprintf(&prompt, "<files>\n%s\n</files>\n", strings.Join(reviewWorkInput.Files, "\n"))

	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range reviewerTools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	conversation := []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(prompt.String()))}
	for turn := 0; turn < maxReviewerTurns; turn++ {
		message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
			Model:     anthropic.BetaManagedAgentsModelClaudeSonnet5_5,
			MaxTokens: int64(4096),
			System:    []anthropic.TextBlockParam{{Text: reviewerSystemPrompt}},
			Messages:  conversation,
			Tools:     anthropicTools,
		})
		if err != nil {
			return "", fmt.Errorf("reviewer inference failed: %w", err)
		}
		Usage.RecordReviewer(message.Usage)
		conversation = append(conversation, message.ToParam())

		var review strings.Builder
		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				review.WriteString(content.Text)
			case "tool_use":
				toolResults = append(toolResults, runReviewerTool(content.ID, content.Name, content.Input))
			}
		}
		if len(toolResults) == 0 {
			return review.String(), nil
		}
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
	}
	return "", fmt.Errorf("reviewer did not reach a verdict within %d turns", maxReviewerTurns)
}

func runReviewerTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	for _, tool := range reviewerTools {
		if tool.Name == name {
			fmt.Printf("\u001b[95mreviewer\u001b[0m: %s(%s)\n", name, input)
			response, err := tool.Function(input)
			if err != nil {
				return anthropic.NewToolResultBlock(id, err.Error(), true)
			}
			return anthropic.NewToolResultBlock(id, response, false)
		}
	}
	return anthropic.NewToolResultBlock(id, "tool not found", true)
}
