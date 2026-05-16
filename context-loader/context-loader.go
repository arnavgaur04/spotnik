package contextloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"spotnik/utils/structs"
	"time"

	"github.com/google/uuid"
)

// LoadContext loads the context from history.jsonl
func LoadContext(maxTurns int) ([]structs.ChatTurn, error) {
	fmt.Println("Starting loading Context")

	// Read from history.jsonl file
	file, err := os.Open("history.jsonl")
	if err != nil {
		if os.IsNotExist(err) {
			return []structs.ChatTurn{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var allTurns []structs.ChatTurn
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var turn structs.ChatTurn
		if err := json.Unmarshal(scanner.Bytes(), &turn); err != nil {
			return nil, err
		}
		allTurns = append(allTurns, turn)
	}

	fmt.Println("Successfully loaded data.jsonl")

	// If history is shorter than maxTurns, return everything
	if len(allTurns) <= maxTurns {
		return allTurns, nil
	}

	// Slice out only the last N elements (e.g., last 10 elements)
	return allTurns[len(allTurns)-maxTurns:], nil
}

// LogTurn appends a chat message directly to history.jsonl
func LogTurn(role, message string) error {
	// Open the file in append-only mode. Create it if it doesn't exist.
	// 0644 grants read/write permissions to the owner.
	file, err := os.OpenFile("history.jsonl", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	logID := uuid.New().String()
	turn := structs.ChatTurn{
		ID:        logID,
		Role:      role,
		Message:   message,
		Timestamp: time.Now(),
	}

	// json.NewEncoder writes directly to the file and appends the trailing \n automatically
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(turn); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
