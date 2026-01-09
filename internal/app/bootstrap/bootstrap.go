package bootstrap

import (
	"context"
	"core-consumer/config"
	"core-consumer/internal/app/infra/db/postgres"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	catalognotification "core-consumer/internal/catalog_notification"
)

func Bootstrap() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	loggerService := logger.New()

	_, err = postgres.Connect(cfg)
	if err != nil {
		panic(err)
	}

	rabbitConn, err := rabbitmq.Connect(cfg, loggerService)
	if err != nil {
		panic(err)
	}

	rabbitConsumer, err := rabbitmq.NewConsumer(rabbitConn, cfg)
	if err != nil {
		panic(err)
	}

	catalognotification.Init(
		rabbitConsumer,
		loggerService,
	)

	rabbitConsumer.Start(context.TODO())

	select {}
}
