package catalognotification

import (
	"core-consumer/internal/app/gen/query"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	browserstorage "core-consumer/internal/app/storage/browser_storage"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotificationprocess "core-consumer/internal/catalog_notification/job/wb_catalog_notification_process"
	wbcatalognotificationstart "core-consumer/internal/catalog_notification/job/wb_catalog_notification_start"
	wbcatalognotification "core-consumer/internal/catalog_notification/repositories/wb_catalog_notification"
)

func Init(
	rabbitConumer *rabbitmq.Consumer,
	loggerService *logger.Logger,
	q *query.Query,
	producer *rabbitmq.Producer,
	browserStorage *browserstorage.Storage,
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

	rabbitConumer.RegisterHandler(
		constants.JobWbCatalogNotificationProccess,
		wbcatalognotificationprocess.New(
			loggerService,
			wbCatalogNotificationRepo,
			producer,
			browserStorage,
		).Handle,
	)

	return nil
}
