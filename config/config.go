package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port             int    `yaml:"port"`
	Model            string `yaml:"model"`
	ContextLimit     int    `yaml:"context_limit"`
	HistoryFile      string `yaml:"history_file"`
	MaxTurns         int    `yaml:"max_turns"`
	ShutdownTimeout  int    `yaml:"shutdown_timeout"`
}

var Current Config

const configPath = "configs/config.yaml"

func Load() error {
	Current = Config{
		Port:            8080,
		Model:           "gemini-3.1-flash-lite-preview",
		ContextLimit:    10,
		HistoryFile:     "history.jsonl",
		MaxTurns:        100,
		ShutdownTimeout: 5,
	}

	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, &Current); err != nil {
			return fmt.Errorf("failed to parse %s: %w", configPath, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	if v := os.Getenv("SPOTNIK_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			Current.Port = n
		}
	}
	if v := os.Getenv("SPOTNIK_MODEL"); v != "" {
		Current.Model = v
	}
	if v := os.Getenv("SPOTNIK_CONTEXT_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			Current.ContextLimit = n
		}
	}
	if v := os.Getenv("SPOTNIK_HISTORY_FILE"); v != "" {
		Current.HistoryFile = v
	}
	if v := os.Getenv("SPOTNIK_MAX_TURNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			Current.MaxTurns = n
		}
	}
	if v := os.Getenv("SPOTNIK_SHUTDOWN_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			Current.ShutdownTimeout = n
		}
	}

	return nil
}
