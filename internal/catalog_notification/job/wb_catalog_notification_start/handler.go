package wbcatalognotificationstart

import (
	"context"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	"log/slog"
)

type Handler struct {
	loggerService *logger.Logger
}

func New(loggerService *logger.Logger) *Handler {
	return &Handler{
		loggerService: loggerService,
	}
}

func (h *Handler) Handle(ctx context.Context, job *rabbitmq.Job) error {
	h.loggerService.Error("HANDLE", slog.String("name", job.Job))

	return nil
}
