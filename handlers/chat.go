package handlers

import (
    "fmt"
    "encoding/json"
    "net/http"
    "spotnik/llm"
    "spotnik/database"
    "spotnik/memory"
    "spotnik/memory/methods"

    "github.com/google/uuid"
    "google.golang.org/genai"
)

type ChatRequest struct {
    User           string `json:"user"`
    Message        string `json:"message"`
    ContextType    int    `json:"context_type"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the incoming request
    var req ChatRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    fmt.Println("request: %+v", req)

    // Generate Task
    taskID := uuid.New().String()
    err := database.CreateTask(taskID)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "ERROR creating task", http.StatusInternalServerError)
        return
    }

    // Save the user's message
    msgID := uuid.New().String()
    err = database.SaveMessage(msgID, taskID, req.Message, req.User)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "ERROR saving message", http.StatusInternalServerError)
        return
    }

    // Load conversation history from DB
    history, err := memory.LoadContext(req.ContextType, req.User, req.Message)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }

    // Build contents with history
    contents := []*genai.Content{}
    for _, h := range history {
        userText := h.Content
        modelText := h.Output

        // Only add to history if both content and output exist
        if h.Content == "" || h.Output == "" {
            continue  // skip incomplete messages
    	}
        contents = append(contents,
            &genai.Content{
                Role:  "user",
                Parts: []*genai.Part{{Text: userText}},
            },
            &genai.Content{
                Role:  "model",
                Parts: []*genai.Part{{Text: modelText}},
            },
        )
    }

    // Add current message
    currentMsg := req.Message
    contents = append(contents, &genai.Content{
        Role:  "user",
        Parts: []*genai.Part{{Text: currentMsg}},
    })

    for i, c := range contents {
    	for _, p := range c.Parts {
            fmt.Printf("Text: %s\n", p.Text)
    	}
    }

    // Call Gemini
    reply, err := llm.CallGemini(contents)
    if err != nil {
        http.Error(w, "LLM error", http.StatusInternalServerError)
        return
    }

    // Save Gemini's reply
    err = database.UpdateMessage(msgID, reply)
    if err != nil {
        http.Error(w, "ERROR saving reply", http.StatusInternalServerError)
        return
    }

    embedding, err := methods.GetEmbedding(reply)
    if err != nil {
        fmt.Println(err)
	http.Error(w, "Embedding Error", http.StatusInternalServerError)
        return
    }

    err = methods.SaveEmbedding(msgID, embedding)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Embedding Save Error", http.StatusInternalServerError)
        return
    }

    err = database.UpdateTaskStatus(taskID, "completed")
    if err != nil {
        fmt.Println(err)
        http.Error(w, "ERROR updating task", http.StatusInternalServerError)
        return
    }

    // Send reply back
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"reply": reply})
}
