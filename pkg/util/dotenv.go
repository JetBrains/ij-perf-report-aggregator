package util

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
	"strings"
)

// LoadDotEnv reads KEY=VALUE pairs from a .env file in the working directory into the process environment.
// Variables that are already set are not overridden, so the file only supplies local development defaults.
// A missing file is not an error.
func LoadDotEnv() {
	f, err := os.Open(".env")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Warn("cannot read .env", "error", err)
		}
		return
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			slog.Warn("cannot set env from .env", "key", key, "error", err)
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Warn("error reading .env", "error", err)
	}
}
