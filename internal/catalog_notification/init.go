package catalognotification

import (
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	"core-consumer/internal/catalog_notification/constants"
	wbcatalognotificationstart "core-consumer/internal/catalog_notification/job/wb_catalog_notification_start"
)

func Init(
	rabbitConumer *rabbitmq.Consumer,
	loggerService *logger.Logger,
) error {
	rabbitConumer.RegisterHandler(
		constants.JobWbCatalogNotificationStart,
		wbcatalognotificationstart.New(loggerService).Handle,
	)

	return nil
}
