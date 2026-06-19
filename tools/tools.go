package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const maxOutputSize = 50 * 1024

func RunTool(call ToolCall) ToolResult {
	switch call.Name {
	case "list_files":
		path, err := call.StringArg("path")
		if err != nil {
			return errorResult(call, err.Error())
		}
		return execResult(call, executeLS(path))

	case "grep_repo":
		query, err := call.StringArg("query")
		if err != nil {
			return errorResult(call, err.Error())
		}
		return execResult(call, executeGrep(query))

	case "cat_file":
		path, err := call.StringArg("path")
		if err != nil {
			return errorResult(call, err.Error())
		}
		return execResult(call, executeCat(path))

	case "write_file":
		path, err := call.StringArg("path")
		if err != nil {
			return errorResult(call, err.Error())
		}
		content, err := call.StringArg("content")
		if err != nil {
			return errorResult(call, err.Error())
		}
		return execResult(call, executeWrite(path, content))

	case "bash":
		command, err := call.StringArg("command")
		if err != nil {
			return errorResult(call, err.Error())
		}
		return execResult(call, executeBash(command))

	default:
		return errorResult(call, fmt.Sprintf("unknown tool: %s", call.Name))
	}
}

func execResult(call ToolCall, output string) ToolResult {
	return ToolResult{ID: call.ID, Name: call.Name, Output: truncate(output)}
}

func errorResult(call ToolCall, msg string) ToolResult {
	return ToolResult{ID: call.ID, Name: call.Name, Error: msg}
}

func truncate(s string) string {
	if len(s) > maxOutputSize {
		return s[:maxOutputSize] + fmt.Sprintf("\n... (truncated, %d bytes total)", len(s))
	}
	return s
}

func executeLS(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}
	var b strings.Builder
	for _, f := range files {
		if f.IsDir() {
			b.WriteString("[DIR]  ")
		} else {
			b.WriteString("[FILE] ")
		}
		b.WriteString(f.Name())
		b.WriteByte('\n')
	}
	return b.String()
}

func executeGrep(query string) string {
	cmd := exec.Command("grep", "-rnE", query, ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "No matches found."
	}
	return string(output)
}

func executeCat(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file %s: %v", path, err)
	}
	return string(content)
}

func executeWrite(path, content string) string {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Sprintf("Error writing to file %s: %v", path, err)
	}
	return fmt.Sprintf("Successfully updated %s", path)
}

func executeBash(command string) string {
	dangerous := []string{"rm -rf /", "mkfs", "dd if=/dev/zero", ":(){:|:&};:"}
	for _, d := range dangerous {
		if strings.Contains(command, d) {
			return "ERROR: Blocked potentially destructive command"
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("ERROR: %v\nOUTPUT: %s", err, string(output))
	}
	return string(output)
}
