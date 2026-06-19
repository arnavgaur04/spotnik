package tools

import "google.golang.org/genai"

func GetGeminiTools() []*genai.Tool {
	defs := GetToolDefs()
	decls := make([]*genai.FunctionDeclaration, len(defs))
	for i, d := range defs {
		decls[i] = GeminiToolDef(d)
	}
	return []*genai.Tool{{
		FunctionDeclarations: decls,
	}}
}

func GetGeminiConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		Tools: GetGeminiTools(),
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode: genai.FunctionCallingConfigModeAuto,
			},
		},
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{
				Text: systemPrompt,
			}},
		},
	}
}
