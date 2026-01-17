package catalognotification

import (
	"core-consumer/internal/app/gen/query"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	browserstorage "core-consumer/internal/app/storage/browser_storage"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotificationprocess "core-consumer/internal/catalog_notification/job/wb_catalog_notification_process"
	wbcatalognotificationstart "core-consumer/internal/catalog_notification/job/wb_catalog_notification_start"
	wbstocknotificationprocess "core-consumer/internal/catalog_notification/job/wb_stock_notification_process"
	wbstocknotificationstart "core-consumer/internal/catalog_notification/job/wb_stock_notification_start"
	wbcatalognotification "core-consumer/internal/catalog_notification/repositories/wb_catalog_notification"
	wbproduct "core-consumer/internal/catalog_notification/repositories/wb_product"
	wbstocknotification "core-consumer/internal/catalog_notification/repositories/wb_stock_notification"
	"core-consumer/internal/stealth"
	telegrambot "core-consumer/internal/telegram_bot"

	"gorm.io/gorm"
)

func Init(
	rabbitConumer *rabbitmq.Consumer,
	loggerService *logger.Logger,
	q *query.Query,
	db *gorm.DB,
	producer *rabbitmq.Producer,
	browserStorage *browserstorage.Storage,
	stealthModule *stealth.Module,
	telgramBotModule *telegrambot.Module,
) error {
	wbCatalogNotificationRepo := wbcatalognotification.New(q)
	rabbitConumer.RegisterHandler(
		constants.JobWbCatalogNotificationStart,
		wbcatalognotificationstart.New(
			loggerService,
			wbCatalogNotificationRepo,
			producer,
		).Handle,
	)

	productsRepo := wbproduct.New(db, q)

	rabbitConumer.RegisterHandler(
		constants.JobWbCatalogNotificationProccess,
		wbcatalognotificationprocess.New(
			loggerService,
			wbCatalogNotificationRepo,
			producer,
			browserStorage,
			stealthModule.ProxyRepo,
			stealthModule.UserAgentRepo,
			productsRepo,
			telgramBotModule.TgBot,
		).Handle,
	)

	wbStockNotificationRepo := wbstocknotification.New(q, db)

	rabbitConumer.RegisterHandler(
		constants.JobWbStockNotificationStart,
		wbstocknotificationstart.New(
			loggerService,
			wbStockNotificationRepo,
			producer,
		).Handle,
	)

	rabbitConumer.RegisterHandler(
		constants.JobWbStockNotificationProccess,
		wbstocknotificationprocess.New(
			loggerService,
			wbStockNotificationRepo,
			producer,
			stealthModule.ProxyRepo,
			stealthModule.UserAgentRepo,
			telgramBotModule.TgBot,
		).Handle,
	)

	return nil
}
