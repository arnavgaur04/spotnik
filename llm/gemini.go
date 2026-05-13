package llm

import (
	"context"
	"encoding/json"
	"fmt"

	Tool "spotnik/tools" // Ensure this matches your module name in go.mod

	"google.golang.org/genai"
)

func CallGemini(contents []*genai.Content) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	// 1. Define the full toolset so the AI knows its capabilities
	tools := []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "grep_repo",
					Description: "Search for text in the local codebase",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {Type: genai.TypeString, Description: "Search string"},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "list_files",
					Description: "List files in the project directory to see structure",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path": {Type: genai.TypeString, Description: "Path to list (use '.' for root)"},
						},
					},
				},
				{
					Name:        "cat_file",
					Description: "Read the full content of a specific file",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path": {Type: genai.TypeString, Description: "Relative path to the file"},
						},
						Required: []string{"path"},
					},
				},
				{
					Name:        "write_file",
					Description: "Overwrite a file with new content. Use this to apply edits.",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path":    {Type: genai.TypeString, Description: "Path to the file"},
							"content": {Type: genai.TypeString, Description: "The full new content of the file"},
						},
						Required: []string{"path", "content"},
					},
				},
			},
		},
	}

	config := &genai.GenerateContentConfig{
		Tools: tools,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{
				Text: "You are an autonomous CLI coding assistant. " +
					"When asked to explain a codebase, do not ask for permission. " +
					"Immediately call list_files('.') to see the structure, " +
					"then use cat_file or grep_repo to understand the logic.",
			}},
		},
	}

	// 2. The Agentic Loop (Thinking -> Acting -> Observing)
	for i := 0; i < 10; i++ { // Increased turns to allow deeper exploration
		result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", contents, config)
		if err != nil {
			// Return error instead of log.Fatal to keep server alive
			return "", fmt.Errorf("generate content error: %w", err)
		}

		// Add the model's response (the 'Thought' or 'Call') to history
		modelResponse := result.Candidates[0].Content
		contents = append(contents, modelResponse)

		hasToolCall := false

		// 3. Process every part of the model's response
		for _, part := range modelResponse.Parts {
			// CRITICAL: Skip parts that are just text (Prevents nil pointer panic)
			if part.FunctionCall == nil {
				continue
			}

			hasToolCall = true
			fmt.Printf("Turn %d: Executing tool %s\n", i+1, part.FunctionCall.Name)

			// 4. Run the actual local command from your tools package
			toolOutput := Tool.RunLocalCommand(part.FunctionCall.Name, part.FunctionCall.Args)

			// 5. Append the observation (Tool Result) to the history
			contents = append(contents, &genai.Content{
				Role: "tool",
				Parts: []*genai.Part{{
					FunctionResponse: &genai.FunctionResponse{
						ID:       part.FunctionCall.ID, // Links response to the specific call
						Name:     part.FunctionCall.Name,
						Response: map[string]any{"result": toolOutput},
					},
				}},
			})
		}

		// Optional: Log the raw response for debugging
		rawJSON, _ := json.MarshalIndent(result, "", "  ")
		_ = rawJSON // Use it or discard it

		// If the AI didn't call any tools, it has provided its final text answer
		if !hasToolCall {
			return result.Text(), nil
		}
	}

	return "The agent reached the maximum number of turns without a final answer.", nil
}
