package tools

import (
	"fmt"
	"os"
	"os/exec"
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
