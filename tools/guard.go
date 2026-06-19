package tools

import (
	"fmt"
	"strings"
)

type Risk string

const (
	Safe   Risk = "safe"
	Prompt Risk = "prompt"
	Block  Risk = "block"
)

type GuardResult struct {
	Risk    Risk
	Message string
}

func Check(call ToolCall, history string) GuardResult {
	switch call.Name {
	case "list_files", "grep_repo", "cat_file":
		return GuardResult{Risk: Safe}

	case "write_file":
		path, err := call.StringArg("path")
		if err != nil {
			return GuardResult{Risk: Block, Message: err.Error()}
		}
		if strings.Contains(path, "..") {
			return GuardResult{
				Risk:    Block,
				Message: fmt.Sprintf("write_file: path %q escapes the project directory", path),
			}
		}
		return GuardResult{Risk: Safe}

	case "bash":
		command, err := call.StringArg("command")
		if err != nil {
			return GuardResult{Risk: Block, Message: err.Error()}
		}

		lower := strings.TrimSpace(strings.ToLower(command))

		if isDestructive(lower) {
			return GuardResult{
				Risk:    Block,
				Message: fmt.Sprintf("bash: command %q is blocked for security", command),
			}
		}

		if isRisky(lower) {
			return GuardResult{
				Risk:    Prompt,
				Message: fmt.Sprintf("PERMISSION_REQUIRED: This will run: %s\nReply with 'yes, proceed' to allow.", command),
			}
		}

		return GuardResult{Risk: Safe}

	default:
		return GuardResult{Risk: Block, Message: fmt.Sprintf("unknown tool: %s", call.Name)}
	}
}

func isDestructive(cmd string) bool {
	patterns := []string{
		"rm -rf /", "rm -rf ~", "rm -rf .",
		"mkfs", "dd if=/dev/zero", "dd if=/dev/random",
		":(){:|:&};:", "> /dev/sda",
		"chmod 000", "chown -r",
	}
	for _, p := range patterns {
		if strings.Contains(cmd, p) {
			return true
		}
	}
	return false
}

func isRisky(cmd string) bool {
	riskyPrefixes := []string{
		"git push", "git commit", "git merge", "git rebase", "git reset",
		"gh pr create", "gh pr merge", "gh repo create",
		"curl", "wget", "ssh", "scp",
		"rm ", "mv ", "sudo",
		"chmod", "chown",
		"> ", ">> ",
		"docker", "kubectl",
		"npm publish", "pip install", "go install",
		"kill", "pkill",
		"shutdown", "reboot",
	}
	for _, p := range riskyPrefixes {
		if strings.HasPrefix(cmd, p) {
			return true
		}
	}
	return false
}
