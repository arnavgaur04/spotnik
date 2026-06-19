package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func RunLocalCommand(name string, args map[string]any) string {
	switch name {
	case "list_files":
		// Run 'ls' or use Go's os.ReadDir
		return executeLS(args["path"].(string))

	case "grep_repo":
		// Run the actual grep command
		return executeGrep(args["query"].(string))

	case "cat_file":
		// Read the file from disk
		return executeCat(args["path"].(string))

	case "write_file":
		path, _ := args["path"].(string)
		content, _ := args["content"].(string)
		if path == "" || content == "" {
			return "Error: path and content are required for write_file"
		}
		return executeWrite(path, content)

	case "bash":
		command, _ := args["command"].(string)
		if command == "" {
			return "Error: command is required"
		}
		return executeBash(command)

	default:
		return "Error: Tool not found"
	}
}

func executeLS(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	var output string
	for _, file := range files {
		if file.IsDir() {
			output += fmt.Sprintf("[DIR]  %s\n", file.Name())
		} else {
			output += fmt.Sprintf("[FILE] %s\n", file.Name())
		}
	}
	return output
}

func executeGrep(query string) string {
	// Runs: grep -rnE "query" .
	// -r: recursive, -n: line numbers, -E: regex
	cmd := exec.Command("grep", "-rnE", query, ".")

	output, err := cmd.CombinedOutput()
	if err != nil {
		// grep returns an error code if no matches are found
		return "No matches found or search error."
	}

	return string(output)
}

func executeCat(path string) string {
	// Read the actual content of the file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file %s: %v", path, err)
	}

	// Return the string content to Gemini
	return string(content)
}

func executeWrite(path string, content string) string {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Error writing to file %s: %v", path, err)
	}
	return fmt.Sprintf("Successfully updated %s", path)
}

func executeBash(command string) string {
	// Basic guard against accidental destruction
	dangerous := []string{"rm -rf /", "mkfs", "dd if=/dev/zero", ":(){:|:&};:"}
	for _, d := range dangerous {
		if strings.Contains(command, d) {
			return "ERROR: Blocked potentially destructive command"
		}
	}

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = "." // lock to project directory

	// Timeout so a bad command doesn't hang forever
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, "bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("ERROR: %v\nOUTPUT: %s", err, string(output))
	}
	return string(output)
}
