package wbcatalognotificationprocess

import (
	"context"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	browserstorage "core-consumer/internal/app/storage/browser_storage"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotification "core-consumer/internal/catalog_notification/parser/wb_catalog_notification"
	wbcatalognotificationRepo "core-consumer/internal/catalog_notification/repositories/wb_catalog_notification"
	"fmt"
	"log/slog"
	"time"
)

type Handler struct {
	loggerService             *logger.Logger
	wbCatalogNotificationRepo *wbcatalognotificationRepo.Repository
	producer                  *rabbitmq.Producer
	browserStorage            *browserstorage.Storage
}

func New(
	loggerService *logger.Logger,
	wbCatalogNotificationRepo *wbcatalognotificationRepo.Repository,
	producer *rabbitmq.Producer,
	browserStorage *browserstorage.Storage,
) *Handler {
	return &Handler{
		loggerService:             loggerService,
		wbCatalogNotificationRepo: wbCatalogNotificationRepo,
		producer:                  producer,
		browserStorage:            browserStorage,
	}
}

func (h *Handler) Handle(ctx context.Context, job *rabbitmq.Job) error {
	h.loggerService.Info("process", slog.String("name", job.Job))
	id, ok := job.Data["id"].(float64)
	if !ok {
		return fmt.Errorf("id not found")
	}

	notification, err := h.wbCatalogNotificationRepo.FindByID(int64(id))
	if err != nil {
		return fmt.Errorf("notification not found")
	}

	if notification.Status != constants.WbCatalogNotificationStatusInProgress {
		h.loggerService.Info("stop", slog.Float64("id", id))

		if err := h.browserStorage.Remove(notification.ID); err != nil {
			h.loggerService.Error(
				"cannot remove browser",
				slog.Float64("id", id),
				slog.String("error", err.Error()),
			)
		}

		return nil
	}

	h.loggerService.Info("slep", slog.Any("sec", notification.Interval))

	parser := wbcatalognotification.New(h.browserStorage)
	products, err := parser.Parse(
		wbcatalognotification.ParseParams{
			NotificationID: notification.ID,
			URL:            "https://www.wildberries.ru/catalog/0/search.aspx?page=1&sort=priceup&search=playstation+5+%D1%81+%D0%B4%D0%B8%D1%81%D0%BA%D0%BE%D0%B2%D0%BE%D0%B4%D0%BE%D0%BC&priceU=3000000%3B250125000&targeturl=ST",
			Proxy:          "93.157.251.104:25029@600523794:804967867",
			UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
			MaxPrice:       50000,
		},
	)
	if err != nil {
		return err
	}

	h.loggerService.Info("success", slog.Any("result", fmt.Sprintf("%v", products)))

	time.Sleep(time.Duration(notification.Interval) * time.Second)

	h.producer.PublishJob(ctx, &rabbitmq.Job{
		Job: constants.JobWbCatalogNotificationProccess,
		Data: map[string]any{
			"id": notification.ID,
		},
	})

	return nil
}
