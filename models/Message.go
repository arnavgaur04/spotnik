package models

type Message struct {
	Role    string         `json:"role"` // "user" or "model"
	Content []ContentBlock `json:"content"`
}
