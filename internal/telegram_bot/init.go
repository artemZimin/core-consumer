package telegrambot

import (
	"core-consumer/config"
	"core-consumer/internal/app/gen/query"
	"core-consumer/internal/telegram_bot/manager/bot"
	telegrambotusers "core-consumer/internal/telegram_bot/repository/telegram_bot_users"
)

type Module struct {
	TgBot *bot.Manager
}

func Init(
	cfg *config.Config,
	q *query.Query,
) (*Module, error) {
	tgBot, err := bot.Load(cfg, telegrambotusers.New(q))
	if err != nil {
		return nil, err
	}

	return &Module{
		TgBot: tgBot,
	}, nil
}
