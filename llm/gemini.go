package llm

import (
	"context"
	"fmt"
	"strings"

	"spotnik/config"
	contextloader "spotnik/context-loader"
	"spotnik/models"
	"spotnik/tools"

	"google.golang.org/genai"
)

func CallGemini(contents []*genai.Content) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	toolConfig := tools.GetGeminiConfig()

	for turn := 0; turn < config.Current.MaxTurns; turn++ {
		result, err := client.Models.GenerateContent(ctx, config.Current.Model, contents, toolConfig)
		if err != nil {
			return "", fmt.Errorf("turn %d: generate content error: %w", turn+1, err)
		}

		modelResponse := result.Candidates[0].Content
		contents = append(contents, modelResponse)

		toolCalls := tools.GeminiCalls(modelResponse)
		if len(toolCalls) == 0 {
			return result.Text(), nil
		}

		modelText := ""
		var toolCallBlocks []models.ContentBlock
		for _, call := range toolCalls {
			toolCallBlocks = append(toolCallBlocks, models.ContentBlock{
				Type:      "tool_use",
				Name:      call.Name,
				Input:     call.Args,
				ToolUseID: call.ID,
			})
		}
		if len(modelResponse.Parts) > 0 && modelResponse.Parts[0].Text != "" {
			modelText = modelResponse.Parts[0].Text
			fmt.Printf("THOUGHT: %s\n", modelText)
		}

		if err := contextloader.LogModelResponse(modelText, toolCallBlocks); err != nil {
			fmt.Printf("WARNING: failed to log model turn: %v\n", err)
		}

		historyText := contentsToString(contents)

		var results []tools.ToolResult
		for _, call := range toolCalls {
			guard := tools.Check(call, historyText)
			var res tools.ToolResult

			switch guard.Risk {
			case tools.Block:
				res = tools.ToolResult{ID: call.ID, Name: call.Name, Error: guard.Message}
				fmt.Printf("TOOL BLOCKED [%s]: %s\n", call.Name, guard.Message)

			case tools.Prompt:
				if containsApproval(historyText) {
					res = tools.RunTool(call)
					if res.IsError() {
						fmt.Printf("TOOL ERROR [%s]: %s\n", call.Name, res.Error)
					} else {
						fmt.Printf("TOOL RESULT [%s]: %d bytes\n", call.Name, len(res.Output))
					}
				} else {
					res = tools.ToolResult{ID: call.ID, Name: call.Name, Output: guard.Message}
					fmt.Printf("TOOL PROMPT [%s]: %s\n", call.Name, guard.Message)
				}

			default:
				res = tools.RunTool(call)
				if res.IsError() {
					fmt.Printf("TOOL ERROR [%s]: %s\n", call.Name, res.Error)
				} else {
					fmt.Printf("TOOL RESULT [%s]: %d bytes\n", call.Name, len(res.Output))
				}
			}

			if err := contextloader.LogToolResult(call.ID, call.Name, res.String()); err != nil {
				fmt.Printf("WARNING: failed to log tool result: %v\n", err)
			}

			results = append(results, res)
		}

		contents = append(contents, &genai.Content{
			Role:  "user",
			Parts: tools.GeminiResultParts(results),
		})
	}

	return "Max turns reached.", nil
}

func containsApproval(text string) bool {
	approvals := []string{
		"yes, proceed", "yes, do it", "approved", "go ahead",
		"i approve", "i approve, proceed", "yeah go ahead",
	}
	lower := strings.ToLower(text)
	for _, a := range approvals {
		if strings.Contains(lower, a) {
			return true
		}
	}
	return false
}

func contentsToString(contents []*genai.Content) string {
	var b strings.Builder
	for _, c := range contents {
		b.WriteString(c.Role)
		b.WriteString(": ")
		for _, p := range c.Parts {
			if p.Text != "" {
				b.WriteString(p.Text)
			}
			if p.FunctionCall != nil {
				b.WriteString("[tool_call: ")
				b.WriteString(p.FunctionCall.Name)
				b.WriteString("]")
			}
			if p.FunctionResponse != nil {
				b.WriteString("[tool_result]")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
