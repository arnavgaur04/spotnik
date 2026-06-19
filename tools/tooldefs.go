package tools

var systemPrompt = `You are an autonomous CLI coding assistant with access to file system tools.

RULES (mandatory):
- NEVER guess file contents — always use cat_file to read them.
- NEVER assume folder structure — always use list_files first.
- NEVER answer questions about code without grounding your answer in tool results.
- If a tool returns an error, try to fix the input and retry once before giving up.

PERMISSION FLOW:
- If a tool returns "PERMISSION_REQUIRED: ...", this means the action needs user approval.
- Ask the user directly: "This will <describe action>. Should I proceed?"
- Wait for their response. If they approve (e.g., "yes, proceed"), retry the tool call.
- If they say no, explain what was blocked and suggest a safer alternative.

FORMAT for every turn:
THOUGHT: <what you know, what you need to find out>
ACTION: <call the appropriate tool>

Wait for the tool result before proceeding.`

func GetToolDefs() []ToolDef {
	return []ToolDef{
		{
			Name:        "list_files",
			Description: "List files in the project directory to see structure. Never guess the repo structure.",
			Parameters: map[string]ParameterDef{
				"path": {Type: "string", Description: "Path to list (use '.' for root)"},
			},
		},
		{
			Name:        "grep_repo",
			Description: "Search for text in the local codebase using regex. Use this to find specific words, phrases, or patterns inside text files.",
			Parameters: map[string]ParameterDef{
				"query": {Type: "string", Description: "Regex search pattern"},
			},
			Required: []string{"query"},
		},
		{
			Name:        "cat_file",
			Description: "Read the full content of a specific file. Use this to understand file contents.",
			Parameters: map[string]ParameterDef{
				"path": {Type: "string", Description: "Relative path to the file"},
			},
			Required: []string{"path"},
		},
		{
			Name:        "write_file",
			Description: "Overwrite a file with new content. Use this to apply edits or create new files.",
			Parameters: map[string]ParameterDef{
				"path":    {Type: "string", Description: "Path to the file"},
				"content": {Type: "string", Description: "The full new content of the file"},
			},
			Required: []string{"path", "content"},
		},
		{
			Name:        "bash",
			Description: "Execute any bash command and get the output. Use this for anything not covered by other tools — running tests, git operations, build commands, or any CLI tool.",
			Parameters: map[string]ParameterDef{
				"command": {Type: "string", Description: "The bash command to run"},
			},
			Required: []string{"command"},
		},
	}
}
