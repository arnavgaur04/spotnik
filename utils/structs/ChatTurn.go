package structs

import (
	"time"
)

// Define the structure of your data
type ChatTurn struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  string    `json:"metadata"`
}
