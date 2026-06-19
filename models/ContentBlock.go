package models

type ContentBlock struct {
	Type      string         `json:"type"` // "text", "tool_use", "tool_result"
	Text      string         `json:"text,omitempty"`
	Name      string         `json:"name,omitempty"`        // tool name
	Input     map[string]any `json:"input,omitempty"`       // tool args
	ToolUseID string         `json:"tool_use_id,omitempty"` // links call to result
	Content   string         `json:"content,omitempty"`     // for tool_result output
}
