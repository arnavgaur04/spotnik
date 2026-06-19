package tools

import "google.golang.org/genai"

func GetGeminiTools() []*genai.Tool {
	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				// {
				// 	Name:        "grep_repo",
				// 	Description: "Search for text in the local codebase/folder. Use this whenver you want to search for specific words, phrases, or patterns inside text files.",
				// 	Parameters: &genai.Schema{
				// 		Type: genai.TypeObject,
				// 		Properties: map[string]*genai.Schema{
				// 			"query": {Type: genai.TypeString, Description: "Search string"},
				// 		},
				// 		Required: []string{"query"},
				// 	},
				// },
				// {
				// 	Name:        "list_files",
				// 	Description: "List files in the project directory to see structure - Never guess the repo structure.",
				// 	Parameters: &genai.Schema{
				// 		Type: genai.TypeObject,
				// 		Properties: map[string]*genai.Schema{
				// 			"path": {Type: genai.TypeString, Description: "Path to list (use '.' for root)"},
				// 		},
				// 	},
				// },
				// {
				// 	Name:        "cat_file",
				// 	Description: "Read the full content of a specific file. Use this whenever you want to understand or go through the file contents.",
				// 	Parameters: &genai.Schema{
				// 		Type: genai.TypeObject,
				// 		Properties: map[string]*genai.Schema{
				// 			"path": {Type: genai.TypeString, Description: "Relative path to the file"},
				// 		},
				// 		Required: []string{"path"},
				// 	},
				// },
				// {
				// 	Name:        "write_file",
				// 	Description: "Overwrite a file with new content. Use this to apply edits or changes.",
				// 	Parameters: &genai.Schema{
				// 		Type: genai.TypeObject,
				// 		Properties: map[string]*genai.Schema{
				// 			"path":    {Type: genai.TypeString, Description: "Path to the file"},
				// 			"content": {Type: genai.TypeString, Description: "The full new content of the file"},
				// 		},
				// 		Required: []string{"path", "content"},
				// 	},
				// },
				{
					Name:        "bash",
					Description: "Execute any bash command and get the output. Use this for anything not covered by other tools — finding files, checking sizes, running tests, git operations, searching, counting lines, anything.",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"command": {Type: genai.TypeString, Description: "The bash command to run"},
						},
						Required: []string{"command"},
					},
				},
			},
		},
	}

	return tools
}

func GetGeminiConfig() *genai.GenerateContentConfig {
	tools := GetGeminiTools()
	config := &genai.GenerateContentConfig{
		Tools: tools,
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode: genai.FunctionCallingConfigModeAuto,
				// Or use FunctionCallingConfigModeAny to FORCE at least one tool call
			},
		},
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{
				Text: `You are an autonomous CLI coding assistant with access to file system tools.
        
        RULES (mandatory):
        - NEVER guess file contents — always use cat_file to read them.
        - NEVER assume folder structure — always use list_files first.
        - NEVER answer questions about code without grounding your answer in tool results.
        - If a tool returns an error, try to fix the input and retry once before giving up.
        
        FORMAT for every turn:
        THOUGHT: <what you know, what you need to find out>
        ACTION: <call the appropriate tool>
        
        Wait for the tool result before proceeding.`,
			}},
		},
	}

	return config
}
