package bootstrap

import (
	"context"
	"core-consumer/config"
	"core-consumer/internal/app/gen/query"
	"core-consumer/internal/app/infra/db/postgres"
	"core-consumer/internal/app/infra/logger"
	"core-consumer/internal/app/infra/queue/rabbitmq"
	browserstorage "core-consumer/internal/app/storage/browser_storage"
	catalognotification "core-consumer/internal/catalog_notification"
	"core-consumer/internal/stealth"
	telegrambot "core-consumer/internal/telegram_bot"
)

func Bootstrap() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	loggerService := logger.New()

	db, err := postgres.Connect(cfg)
	if err != nil {
		panic(err)
	}

	q := query.Use(db)

	rabbitConn, err := rabbitmq.Connect(cfg, loggerService)
	if err != nil {
		panic(err)
	}

	rabbitProducer, err := rabbitmq.NewProducer(rabbitConn, cfg)

	rabbitConsumer, err := rabbitmq.NewConsumer(rabbitConn, cfg)
	if err != nil {
		panic(err)
	}

	browserStorage := browserstorage.New(cfg)

	stealthModule := stealth.Init(db, q)
	telegramBotModule, err := telegrambot.Init(
		cfg,
		q,
	)
	if err != nil {
		panic(err)
	}

	catalognotification.Init(
		rabbitConsumer,
		loggerService,
		q,
		db,
		rabbitProducer,
		browserStorage,
		stealthModule,
		telegramBotModule,
	)

	rabbitConsumer.Start(context.TODO())

	select {}
}
