package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var CreateFileDefinition = ToolDefinition{
	Name: "create_file",
	Description: "Create a file at the given relative path. Supply ./ at the beginning of the path. Truncates the file if the file already exists",
	InputSchema: CreateFileInputSchema,
	Function: CreateFile,
}

type CreateFileInput struct {
	Path	string `json:"path" jsonschema_description:"The relative path of the file to be created"`
}

var CreateFileInputSchema = GenerateSchema[CreateFileInput]()

func CreateFile(input json.RawMessage) (string, error) {
	createFileInput := CreateFileInput{}
	err := json.Unmarshal(input, &createFileInput)

	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(createFileInput.Path), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	if createFileInput.Path == "" {
		return "", fmt.Errorf("Invalid path")
	}

	file, err := os.Create(createFileInput.Path)
	if err != nil {
		return "", err
	}

	file.Close()

	return fmt.Sprintf("File successfully created at %s", createFileInput.Path), nil
}
