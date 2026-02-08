package wbstocknotificationprocess

import (
	"context"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotification "core-consumer/internal/catalog_notification/parser/wb_catalog_notification"
	wbproductprice "core-consumer/internal/catalog_notification/repositories/wb_product_price"
	wbstocknotification "core-consumer/internal/catalog_notification/repositories/wb_stock_notification"
	"core-consumer/internal/stealth/repository/proxy"
	useragent "core-consumer/internal/stealth/repository/user_agent"
	"core-consumer/internal/telegram_bot/manager/bot"
	"fmt"
	"log/slog"
	"time"
)

type Handler struct {
	loggerService           *logger.Logger
	wbStockNotificationRepo *wbstocknotification.Repository
	producer                *rabbitmq.Producer
	proxyRepo               *proxy.Repository
	userAgentRepo           *useragent.Repository
	tgBot                   *bot.Manager
	wbProductPriceRepo      *wbproductprice.Repository
}

func New(
	loggerService *logger.Logger,
	wbStockNotificationRepo *wbstocknotification.Repository,
	producer *rabbitmq.Producer,
	proxyRepo *proxy.Repository,
	userAgentRepo *useragent.Repository,
	tgBot *bot.Manager,
	wbProductPriceRepo *wbproductprice.Repository,
) *Handler {
	return &Handler{
		loggerService:           loggerService,
		wbStockNotificationRepo: wbStockNotificationRepo,
		producer:                producer,
		proxyRepo:               proxyRepo,
		userAgentRepo:           userAgentRepo,
		tgBot:                   tgBot,
		wbProductPriceRepo:      wbProductPriceRepo,
	}
}

func (h *Handler) Handle(ctx context.Context, job *rabbitmq.Job) error {
	h.loggerService.Info("process", slog.String("name", job.Job))
	id, ok := job.Data["id"].(float64)
	if !ok {
		return fmt.Errorf("id not found")
	}

	countNotifications, err := h.wbStockNotificationRepo.CountInStatus(
		constants.WbCatalogNotificationStatusInProgress,
	)
	if err != nil {
		return fmt.Errorf("fail countNotifications")
	}

	notification, err := h.wbStockNotificationRepo.FindByID(int64(id))
	if err != nil {
		return fmt.Errorf("notification not found")
	}

	if notification.Status != constants.WbCatalogNotificationStatusInProgress {
		h.loggerService.Info("stop", slog.Float64("id", id))

		return nil
	}

	if notification.Cookie == nil {
		return fmt.Errorf("cookie not set for stock notification %d", notification.ID)
	}

	h.loggerService.Info("slep", slog.Any("sec", notification.Interval))

	proxy, err := h.proxyRepo.FindRanbom()
	if err != nil {
		h.loggerService.Error("proxy not found", slog.String("error", err.Error()))
	}

	userAgent, err := h.userAgentRepo.FindRandom()
	if err != nil {
		h.loggerService.Error("user agent not found", slog.String("error", err.Error()))
	}

	parser := wbcatalognotification.New(nil)

	h.loggerService.Info(
		"prepare parse wb catalog",
		slog.Int64("stock_notification", notification.ID),
		slog.Int64("stock_notification_count", countNotifications),
		slog.Int64("proxy", proxy.ID),
		slog.Int64("user_agent", userAgent.ID),
	)

	products, err := parser.ParseStock(
		wbcatalognotification.ParseStockParams{
			URL:       notification.URL,
			Proxy:     proxy.Data,
			UserAgent: userAgent.Data,
			Cookie:    *notification.Cookie,
		},
	)
	if err != nil {
		h.loggerService.Error(
			"parse error",
			slog.String("error", err.Error()),
			slog.Int64("stock_notification", notification.ID),
			slog.Int64("stock_notification_count", countNotifications),
			slog.Int64("proxy", proxy.ID),
			slog.Int64("user_agent", userAgent.ID),
		)

		h.producer.PublishJob(ctx, &rabbitmq.Job{
			Job: constants.JobWbStockNotificationProccess,
			Data: map[string]any{
				"id": notification.ID,
			},
		})

		return nil
	}

	h.loggerService.Info(
		"stock_notification success",
		slog.Int("products_count", len(products)),
		slog.Int64("stock_notification", notification.ID),
		slog.Int64("stock_notification_count", countNotifications),
		slog.Int64("proxy", proxy.ID),
		slog.Int64("user_agent", userAgent.ID),
	)

	for _, product := range products {
		isInStock := product.Quantity > 0

		if notification.MaxPrice != nil && int64(*notification.MaxPrice) < product.Price {
			continue
		}

		if notification.IsInStock == false && isInStock {
			if err := h.tgBot.BroadcastWbStockNotification(
				bot.BroadcastWbStockNotificationParam{
					ProductName:      product.Name,
					ImgURL:           product.Img,
					NotificationName: notification.Name,
					ProductURL:       product.URL,
					Price:            product.Price,
					Quantity:         product.Quantity,
				},
			); err != nil {
				h.loggerService.Error(
					"Failed broadcast wb notification",
					slog.Int64("wb_stock_notification_id", notification.ID),
				)
			}
		}
	}

	var maxPrice *int32

	if notification.WbProductPriceID != nil {
		price, err := h.wbProductPriceRepo.FindByID(*notification.WbProductPriceID)
		if err != nil {
			return err
		}
		currentPrice := int32(float32(price.MaxPrice) * 1.1)

		maxPrice = &currentPrice
	} else {
		maxPrice = notification.MaxPrice
	}

	isInStock := false
	for _, product := range products {
		isInStock = product.Quantity > 0 && (maxPrice == nil || int64(*maxPrice) > product.Price)
		if isInStock {
			break
		}
	}
	err = h.wbStockNotificationRepo.UpdateIsInStock(notification.ID, isInStock)
	if err != nil {
		h.loggerService.Error(
			"Failed wbStockNotificationRepo.UpdateIsInStock",
			slog.Int64("wb_stock_notification_id", notification.ID),
		)
	}

	time.Sleep(time.Duration(notification.Interval) * time.Second)

	h.producer.PublishJob(ctx, &rabbitmq.Job{
		Job: constants.JobWbStockNotificationProccess,
		Data: map[string]any{
			"id": notification.ID,
		},
	})

	return nil
}
