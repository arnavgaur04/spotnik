package llm

import (
	"context"
	"fmt"
	"strings"

	"spotnik/config"
	contextloader "spotnik/context-loader"
	Tool "spotnik/tools"
	"spotnik/models"

	"google.golang.org/genai"
)

func CallGemini(contents []*genai.Content) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	// 1. Define the full toolset so the AI knows its capabilities
	toolConfig := Tool.GetGeminiConfig()

	// 2. The Agentic Loop (Thinking -> Acting -> Observing)
	for turn := 0; turn < config.Current.MaxTurns; turn++ {
		result, err := client.Models.GenerateContent(ctx, config.Current.Model, contents, toolConfig)
		if err != nil {
			return "", fmt.Errorf("turn %d: generate content error: %w", turn+1, err)
		}

		modelResponse := result.Candidates[0].Content
		contents = append(contents, modelResponse)

		var toolCallBlocks []models.ContentBlock
		var modelText string
		hasToolCall := false

		// First pass — collect everything, log model turn FIRST
		for _, part := range modelResponse.Parts {
			if part.Text != "" {
				modelText = part.Text
				fmt.Printf("THOUGHT: %s\n", part.Text)
			}
			if part.FunctionCall != nil {
				hasToolCall = true
				toolCallBlocks = append(toolCallBlocks, models.ContentBlock{
					Type:      "tool_use",
					Name:      part.FunctionCall.Name,
					Input:     part.FunctionCall.Args,
					ToolUseID: part.FunctionCall.ID,
				})
			}
		}

		// ✅ Log model turn BEFORE any tool execution or tool result logging
		if err := contextloader.LogModelResponse(modelText, toolCallBlocks); err != nil {
			fmt.Printf("WARNING: failed to log model turn: %v\n", err)
		}

		// Second pass — execute tools and log results AFTER model turn
		var funcRespParts []*genai.Part
		for _, part := range modelResponse.Parts {
			if part.FunctionCall == nil {
				continue
			}

			toolOutput := Tool.RunLocalCommand(part.FunctionCall.Name, part.FunctionCall.Args)
			if strings.HasPrefix(toolOutput, "ERROR:") {
				fmt.Printf("TOOL ERROR [%s]: %s\n", part.FunctionCall.Name, toolOutput)
			} else {
				fmt.Printf("TOOL RESULT [%s]: %s\n", part.FunctionCall.Name, toolOutput)
			}

			if err := contextloader.LogToolResult(part.FunctionCall.ID, part.FunctionCall.Name, toolOutput); err != nil {
				fmt.Printf("WARNING: failed to log tool result: %v\n", err)
			}

			funcRespParts = append(funcRespParts, &genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					ID:       part.FunctionCall.ID,
					Name:     part.FunctionCall.Name,
					Response: map[string]any{"result": toolOutput},
				},
			})
		}
		if len(funcRespParts) > 0 {
			contents = append(contents, &genai.Content{
				Role:  "user",
				Parts: funcRespParts,
			})
		}

		if !hasToolCall {
			return result.Text(), nil
		}
	}

	return "Max turns reached.", nil
}
