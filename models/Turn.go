package models

import "time"

type Turn struct {
	UUID      string    `json:"uuid"`
	Timestamp time.Time `json:"timestamp"`
	Message   Message   `json:"message"`
}
