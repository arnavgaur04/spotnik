package handlers

import (
    "encoding/json"
    "net/http"
    "spotnik/llm"
)

type ChatRequest struct {
    SessionID string `json:"session_id"`  // who is talking
    Message   string `json:"message"`     // what they said
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse the incoming request
    var req ChatRequest
    json.NewDecoder(r.Body).Decode(&req)

    // 2. Load conversation history from DB
    //history, err := database.LoadHistory(req.SessionID)
    //if err != nil {
    //    http.Error(w, "DB error", http.StatusInternalServerError)
    //    return
    //}

    // 3. Save the user's message
    //database.SaveMessage(req.SessionID, "user", req.Message)

    // 4. Call Gemini with history + new message
    reply, err := llm.CallGemini(req.Message)
    if err != nil {
        http.Error(w, "LLM error", http.StatusInternalServerError)
        return
    }

    // 5. Save Gemini's reply
    //database.SaveMessage(req.SessionID, "assistant", reply)

    // 6. Send reply back
    json.NewEncoder(w).Encode(map[string]string{"reply": reply})
}
