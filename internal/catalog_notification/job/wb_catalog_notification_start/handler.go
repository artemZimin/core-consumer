package wbcatalognotificationstart

import (
	"context"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotification "core-consumer/internal/catalog_notification/repositories/wb_catalog_notification"
	"fmt"
	"log/slog"
)

type Handler struct {
	loggerService             *logger.Logger
	wbCatalogNotificationRepo *wbcatalognotification.Repository
	producer                  *rabbitmq.Producer
}

func New(
	loggerService *logger.Logger,
	wbCatalogNotificationRepo *wbcatalognotification.Repository,
	producer *rabbitmq.Producer,
) *Handler {
	return &Handler{
		loggerService:             loggerService,
		wbCatalogNotificationRepo: wbCatalogNotificationRepo,
		producer:                  producer,
	}
}

func (h *Handler) Handle(ctx context.Context, job *rabbitmq.Job) error {
	h.loggerService.Error("HANDLE", slog.String("name", job.Job))
	id, ok := job.Data["id"].(float64)
	if !ok {
		return fmt.Errorf("id not found")
	}

	notification, err := h.wbCatalogNotificationRepo.FindByID(int64(id))
	if err != nil {
		return fmt.Errorf("notification not found")
	}

	h.producer.PublishJob(ctx, &rabbitmq.Job{
		Job: constants.JobWbCatalogNotificationProccess,
		Data: map[string]any{
			"id": notification.ID,
		},
	})

	return nil
}
