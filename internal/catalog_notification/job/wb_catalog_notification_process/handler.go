package wbcatalognotificationprocess

import (
	"context"
	"core-consumer/internal/app/gen/model"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	browserstorage "core-consumer/internal/app/storage/browser_storage"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotification "core-consumer/internal/catalog_notification/parser/wb_catalog_notification"
	wbcatalognotificationRepo "core-consumer/internal/catalog_notification/repositories/wb_catalog_notification"
	wbproduct "core-consumer/internal/catalog_notification/repositories/wb_product"
	"core-consumer/internal/stealth/repository/proxy"
	useragent "core-consumer/internal/stealth/repository/user_agent"
	"core-consumer/internal/telegram_bot/manager/bot"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Handler struct {
	loggerService             *logger.Logger
	wbCatalogNotificationRepo *wbcatalognotificationRepo.Repository
	producer                  *rabbitmq.Producer
	browserStorage            *browserstorage.Storage
	proxyRepo                 *proxy.Repository
	userAgentRepo             *useragent.Repository
	productsRepo              *wbproduct.Repository
	tgBot                     *bot.Manager
}

func New(
	loggerService *logger.Logger,
	wbCatalogNotificationRepo *wbcatalognotificationRepo.Repository,
	producer *rabbitmq.Producer,
	browserStorage *browserstorage.Storage,
	proxyRepo *proxy.Repository,
	userAgentRepo *useragent.Repository,
	productsRepo *wbproduct.Repository,
	tgBot *bot.Manager,
) *Handler {
	return &Handler{
		loggerService:             loggerService,
		wbCatalogNotificationRepo: wbCatalogNotificationRepo,
		producer:                  producer,
		browserStorage:            browserStorage,
		proxyRepo:                 proxyRepo,
		userAgentRepo:             userAgentRepo,
		productsRepo:              productsRepo,
		tgBot:                     tgBot,
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

		return nil
	}

	if notification.Cookie == nil {
		return fmt.Errorf("cookie not set for notification %d", notification.ID)
	}

	h.loggerService.Info("slep", slog.Any("sec", notification.Interval))

	proxy, err := h.proxyRepo.FindByWbCatalogNotification(notification.ID)
	if err != nil {
		h.loggerService.Error("proxy not found", slog.String("error", err.Error()))
	}

	userAgent, err := h.userAgentRepo.FindRandom()
	if err != nil {
		h.loggerService.Error("user agent not found", slog.String("error", err.Error()))
	}

	parser := wbcatalognotification.New(h.browserStorage)

	h.loggerService.Info(
		"prepare parse wb catalog",
		slog.Int64("notification", notification.ID),
		slog.Int64("proxy", proxy.ID),
		slog.Int64("user_agent", userAgent.ID),
	)

	products, err := parser.ParseV2(
		wbcatalognotification.ParseParams{
			NotificationID: notification.ID,
			URL:            notification.URL,
			Proxy:          proxy.Data,
			UserAgent:      userAgent.Data,
			MaxPrice:       int64(notification.MaxPrice),
			Cookie:         *notification.Cookie,
		},
	)
	if err != nil {
		h.loggerService.Error(
			"parse error",
			slog.String("error", err.Error()),
			slog.Int64("notification", notification.ID),
			slog.Int64("proxy", proxy.ID),
			slog.Int64("user_agent", userAgent.ID),
		)

		h.producer.PublishJob(ctx, &rabbitmq.Job{
			Job: constants.JobWbCatalogNotificationProccess,
			Data: map[string]any{
				"id": notification.ID,
			},
		})

		return nil
	}

	h.loggerService.Info(
		"catalog_notification success",
		slog.Int("products_count", len(products)),
		slog.Int64("catalog_notification", notification.ID),
		slog.Int64("proxy", proxy.ID),
		slog.Int64("user_agent", userAgent.ID),
	)

ProductsLoop:
	for _, product := range products {
		if notification.MinPrice != nil && int64(*notification.MinPrice) > product.Price {
			continue
		}

		if int64(notification.MaxPrice) < product.Price {
			continue
		}

		if notification.StopWords != nil {
			stopWords := strings.Split(strings.ToLower(*notification.StopWords), ",")
			name := strings.ToLower(product.Name)

			for _, stopWord := range stopWords {
				if strings.Contains(name, stopWord) {
					continue ProductsLoop
				}
			}
		}

		if notification.PlusWords != nil {
			plusWords := strings.Split(strings.ToLower(*notification.PlusWords), ",")
			name := strings.ToLower(product.Name)

			for _, plusWord := range plusWords {
				if !strings.Contains(name, plusWord) {
					continue ProductsLoop
				}
			}
		}

		_, err := h.productsRepo.FindByUrlAndPriceInCatalogNotification(
			wbproduct.FindByUrlAndPriceInCatalogNotificationParams{
				NotificationID: notification.ID,
				URL:            product.URL,
				Price:          int32(product.Price),
			},
		)
		if err == nil {
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		h.loggerService.Info("product", slog.Any("result", fmt.Sprintf("%v", product)))
		err = h.productsRepo.Create(&model.WbCatalogNotificationProduct{
			Price:                   int32(product.Price),
			URL:                     product.URL,
			Img:                     product.Img,
			WbCatalogNotificationID: notification.ID,
		})
		if err != nil {
			return err
		}

		if err := h.tgBot.BroadcastWbCatalogNotification(
			bot.BroadcastWbCatalogNotificationParam{
				ImgURL:           product.Img,
				NotificationName: notification.Name,
				ProductURL:       product.URL,
				Price:            product.Price,
				Quantity:         product.Quantity,
			},
		); err != nil {
			h.loggerService.Error(
				"Failed broadcast wb notification",
				slog.Int64("wb_catalog_notification_id", notification.ID),
			)
		}
	}

	time.Sleep(time.Duration(notification.Interval) * time.Second)

	h.producer.PublishJob(ctx, &rabbitmq.Job{
		Job: constants.JobWbCatalogNotificationProccess,
		Data: map[string]any{
			"id": notification.ID,
		},
	})

	return nil
}
