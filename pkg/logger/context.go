package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const loggerKey contextKey = "logger"

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts logger from context, returns default if not found
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	// Return default logger if not found in context
	return slog.Default()
}
