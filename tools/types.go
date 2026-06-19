package tools

import "fmt"

type ToolCall struct {
	ID   string
	Name string
	Args map[string]any
}

type ToolResult struct {
	ID     string
	Name   string
	Output string
	Error  string
}

func (r ToolResult) String() string {
	if r.Error != "" {
		return r.Error
	}
	return r.Output
}

func (r ToolResult) IsError() bool {
	return r.Error != ""
}

type ToolDef struct {
	Name        string
	Description string
	Parameters  map[string]ParameterDef
	Required    []string
}

type ParameterDef struct {
	Type        string
	Description string
}

func (c *ToolCall) StringArg(key string) (string, error) {
	v, ok := c.Args[key]
	if !ok {
		return "", fmt.Errorf("tool %s: missing required argument %q", c.Name, key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("tool %s: argument %q must be a string, got %T", c.Name, key, v)
	}
	return s, nil
}
