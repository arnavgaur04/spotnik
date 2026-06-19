package models

type ChatRequest struct {
	User        string `json:"user"`
	Message     string `json:"message"`
	ContextType int    `json:"context_type"`
}
