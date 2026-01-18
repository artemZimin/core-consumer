package logger

import (
	"context"
	"core-consumer/config"
	"io"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
	cfg    *config.Config
}

func New(cfg *config.Config, out io.Writer) *Logger {
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
		cfg: cfg,
	}
}

func (s *Logger) log(level slog.Level, msg string, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("service", s.cfg.AppName))

	s.logger.LogAttrs(
		context.Background(),
		level,
		msg,
		attrs...,
	)
}

func (s *Logger) Error(msg string, attrs ...slog.Attr) {
	s.log(
		slog.LevelError,
		msg,
		attrs...,
	)
}

func (s *Logger) Info(msg string, attrs ...slog.Attr) {
	s.log(
		slog.LevelInfo,
		msg,
		attrs...,
	)
}

func (s *Logger) Debug(msg string, attrs ...slog.Attr) {
	s.log(
		slog.LevelDebug,
		msg,
		attrs...,
	)
}
