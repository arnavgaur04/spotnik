package tools

import (
	"google.golang.org/genai"
)

func GeminiCalls(c *genai.Content) []ToolCall {
	var calls []ToolCall
	for _, part := range c.Parts {
		if part.FunctionCall == nil {
			continue
		}
		calls = append(calls, ToolCall{
			ID:   part.FunctionCall.ID,
			Name: part.FunctionCall.Name,
			Args: part.FunctionCall.Args,
		})
	}
	return calls
}

func GeminiResultParts(results []ToolResult) []*genai.Part {
	parts := make([]*genai.Part, len(results))
	for i, r := range results {
		parts[i] = &genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				ID:   r.ID,
				Name: r.Name,
				Response: map[string]any{
					"result": r.String(),
				},
			},
		}
	}
	return parts
}

func GeminiToolDef(t ToolDef) *genai.FunctionDeclaration {
	props := make(map[string]*genai.Schema)
	for name, p := range t.Parameters {
		props[name] = &genai.Schema{
			Type:        genai.TypeString,
			Description: p.Description,
		}
	}
	return &genai.FunctionDeclaration{
		Name:        t.Name,
		Description: t.Description,
		Parameters: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: props,
			Required:   t.Required,
		},
	}
}
