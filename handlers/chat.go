package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	contextloader "spotnik/context-loader"
	"spotnik/llm"
	"spotnik/utils/structs"

	"google.golang.org/genai"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var req structs.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Creating an entry for user prompt in history.jsonl
	err := contextloader.LogTurn("user", req.Message)
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

	// Build contents with context
	contents := []*genai.Content{}
	for _, c := range context {
		// Only add to context if message exist
		if c.Message == "" {
			continue // skip incomplete messages
		}
		contents = append(contents,
			&genai.Content{
				Role:  c.Role,
				Parts: []*genai.Part{{Text: c.Message}},
			},
		)
	}

	// Call Gemini
	reply, err := llm.CallGemini(contents)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "LLM error", http.StatusInternalServerError)
		return
	}

	// Save Gemini's reply
	err = contextloader.LogTurn("model", reply)
	if err != nil {
		http.Error(w, "ERROR saving reply", http.StatusInternalServerError)
		return
	}

	// Send reply back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"reply": reply})
}
