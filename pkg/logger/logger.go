package logger

import (
	"log/slog"
	"os"
)

// NewLogger creates a structured JSON logger
func NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// NewTestLogger creates a logger for tests (text format)
func NewTestLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
