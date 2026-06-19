package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	contextloader "spotnik/context-loader"
	"spotnik/llm"
	"spotnik/models"

	"google.golang.org/genai"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Creating an entry for user prompt in history.jsonl
	err := contextloader.LogUserMessage(req.Message)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ERROR logging prompt", http.StatusInternalServerError)
		return
	}

	// Load conversation history from history.jsonl
	context, err := contextloader.LoadContext(10)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ERROR loading context", http.StatusInternalServerError)
		return
	}

	var contents []*genai.Content

	// Rebuild contents from history, skipping orphaned function responses
	pendingFC := false
	for _, turn := range context {
		switch turn.Message.Role {
		case "user":
			var textParts []*genai.Part
			var funcRespParts []*genai.Part
			for _, block := range turn.Message.Content {
				if block.Type == "text" {
					textParts = append(textParts, &genai.Part{Text: block.Text})
				} else if block.Type == "tool_result" {
					if pendingFC {
						funcRespParts = append(funcRespParts, &genai.Part{
							FunctionResponse: &genai.FunctionResponse{
								ID:       block.ToolUseID,
								Name:     block.Name,
								Response: map[string]any{"result": block.Content},
							},
						})
					}
				}
			}
			if len(funcRespParts) > 0 {
				contents = append(contents, &genai.Content{
					Role:  "user",
					Parts: funcRespParts,
				})
				pendingFC = false
			}
			if len(textParts) > 0 {
				contents = append(contents, &genai.Content{
					Role:  "user",
					Parts: textParts,
				})
				pendingFC = false
			}

		case "model":
			var parts []*genai.Part
			for _, block := range turn.Message.Content {
				if block.Type == "text" {
					parts = append(parts, &genai.Part{Text: block.Text})
				} else if block.Type == "tool_use" {
					parts = append(parts, &genai.Part{
						FunctionCall: &genai.FunctionCall{
							ID:   block.ToolUseID,
							Name: block.Name,
							Args: block.Input,
						},
					})
				}
			}
			if len(parts) > 0 {
				hasFC := false
				for _, p := range parts {
					if p.FunctionCall != nil {
						hasFC = true
						break
					}
				}
				pendingFC = hasFC
				contents = append(contents, &genai.Content{
					Role:  "model",
					Parts: parts,
				})
			}
		}
	}

	// Call Gemini
	contents = sanitizeContents(contents)
	reply, err := llm.CallGemini(contents)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "LLM error", http.StatusInternalServerError)
		return
	}

	// Send reply back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"reply": reply})
}

func sanitizeContents(contents []*genai.Content) []*genai.Content {
	var sanitized []*genai.Content

	for _, c := range contents {
		if len(sanitized) == 0 {
			if c.Role != "user" {
				continue
			}
			hasText := false
			for _, p := range c.Parts {
				if p.Text != "" {
					hasText = true
					break
				}
			}
			if !hasText {
				continue
			}
			sanitized = append(sanitized, c)
			continue
		}

		last := sanitized[len(sanitized)-1]

		if last.Role == c.Role && c.Role != "user" {
			fmt.Printf("WARNING: skipping duplicate role '%s'\n", c.Role)
			continue
		}

		sanitized = append(sanitized, c)
	}

	return sanitized
}
