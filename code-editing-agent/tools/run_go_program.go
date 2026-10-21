package tools

import (
	"encoding/json"
	"os/exec"
)

var RunProgramDefinition = ToolDefinition {
	Name: 			"run_program",
	Description: 	"Run a program based on the relative filepath on the local filesystem You must yourself append the ./. Use when you want to execute a go program, or capture the output for further validation and/or testing. Options are either \"go\" or \"node\" and NOTHING else",
	InputSchema:	RunProgramSchema ,
	Function: 		RunProgram,
}

type RunProgramInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of the file you want to run"`
	Program string `json:"program" jsonschema_description:"The terminal program you want to run the program with. Options are either "go" or "node". This will be appended infront of the path variable that is supplied to run the program at the path location. The harness supplies either run or nothing as the second argument"`
	Args []string `json:"args" jsonschema_description:"The arguments supplied to the program. Must not interfere with or override or extend the original program invocations like go run or node commands. must only be arguments supplied to the artifact under execution"`
}

var RunProgramSchema = GenerateSchema[RunProgramInput]()

func RunProgram(input json.RawMessage) (string, error) {
	runProgramInput := RunProgramInput{}
	err := json.Unmarshal(input, &runProgramInput)

	if err != nil {
		panic(err)
	}

	path := runProgramInput.Path
	program := runProgramInput.Program
	var args []string

	if(program == "go") {
		args = []string{"run", path}
	} else {
		args = []string{path}
	}

	args = append(args, runProgramInput.Args...)

	cmd := exec.Command(program, args...)

	byte_arr, err := cmd.Output()

	if err != nil {
		return string(""), err
	}

	output := string(byte_arr)

	return output, nil
}
