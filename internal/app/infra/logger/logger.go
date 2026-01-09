package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New() *Logger {
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}
}

func (s *Logger) log(level slog.Level, msg string, attrs ...slog.Attr) {
	s.logger.LogAttrs(
		context.Background(),
		level,
		msg,
		attrs...,
	)
}

func (s *Logger) Error(msg string, attrs ...slog.Attr) {
	s.log(slog.LevelError, msg, attrs...)
}

func (s *Logger) Info(msg string, attrs ...slog.Attr) {
	s.log(slog.LevelInfo, msg, attrs...)
}

func (s *Logger) Debug(msg string, attrs ...slog.Attr) {
	s.log(slog.LevelDebug, msg, attrs...)
}
