package contextloader

import (
	"bufio"
	"encoding/json"
	"os"
	"spotnik/models"
	"time"

	"github.com/google/uuid"
)

const historyFile = "history.jsonl"

func appendTurn(turn models.Turn) error {
	f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(turn)
}

// LogUserMessage - plain text from user
func LogUserMessage(text string) error {
	return appendTurn(models.Turn{
		UUID:      uuid.NewString(),
		Timestamp: time.Now(),
		Message: models.Message{
			Role: "user",
			Content: []models.ContentBlock{
				{Type: "text", Text: text},
			},
		},
	})
}

// LogModelResponse - model text + any tool calls together in one turn
func LogModelResponse(text string, toolCalls []models.ContentBlock) error {
	content := []models.ContentBlock{}
	if text != "" {
		content = append(content, models.ContentBlock{
			Type: "text",
			Text: text,
		})
	}
	content = append(content, toolCalls...)
	return appendTurn(models.Turn{
		UUID:      uuid.NewString(),
		Timestamp: time.Now(),
		Message: models.Message{
			Role:    "model",
			Content: content,
		},
	})
}

// LogToolResult - tool result goes in a "user" role message (same as Claude Code)
func LogToolResult(toolUseID, toolName, result string) error {
	return appendTurn(models.Turn{
		UUID:      uuid.NewString(),
		Timestamp: time.Now(),
		Message: models.Message{
			Role: "user",
			Content: []models.ContentBlock{
				{
					Type:      "tool_result",
					ToolUseID: toolUseID,
					Name:      toolName,
					Content:   result,
				},
			},
		},
	})
}

// LoadContext loads last N turns from history
func LoadContext(limit int) ([]models.Turn, error) {
	f, err := os.Open(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var turns []models.Turn
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var turn models.Turn
		if err := json.Unmarshal(scanner.Bytes(), &turn); err != nil {
			continue
		}
		turns = append(turns, turn)
	}

	// Return last N turns
	if len(turns) > limit {
		turns = turns[len(turns)-limit:]
	}
	return turns, scanner.Err()
}
