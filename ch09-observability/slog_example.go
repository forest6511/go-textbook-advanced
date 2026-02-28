package main

import (
	"log/slog"
	"os"
	"time"
)

// buildDBAttrs builds a slice of slog.Attr for DB operations
func buildDBAttrs(query string, dur time.Duration, err error) []slog.Attr {
	attrs := []slog.Attr{
		slog.String("query", query),
		slog.Duration("duration", dur),
	}
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	return attrs
}

// newMultiLogger creates a logger that writes to multiple handlers (Go 1.26)
func newMultiLogger() *slog.Logger {
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	fileHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	multi := slog.NewMultiHandler(consoleHandler, fileHandler)
	return slog.New(multi)
}

func main() {
	logger := newMultiLogger()

	// GroupAttrs (Go 1.25): create a group Attr from a slice using spread
	dbAttrs := buildDBAttrs("SELECT * FROM users", 12*time.Millisecond, nil)
	logger.Info("db query", slog.GroupAttrs("db", dbAttrs...))

	// Standard structured logging
	logger.Info("server started",
		slog.String("addr", ":8080"),
		slog.Int("pid", os.Getpid()),
	)

	// Dynamic log level
	var level slog.LevelVar
	level.Set(slog.LevelDebug)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &level})
	debugLogger := slog.New(handler)
	debugLogger.Debug("debug message", slog.String("key", "value"))
}
