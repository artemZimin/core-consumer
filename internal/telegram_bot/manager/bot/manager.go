package bot

import (
	"core-consumer/config"
	telegrambotusers "core-consumer/internal/telegram_bot/repository/telegram_bot_users"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Manager struct {
	api       *tgbotapi.BotAPI
	usersRepo *telegrambotusers.Repository
}

func Load(
	cfg *config.Config,
	usersRepo *telegrambotusers.Repository,
) (*Manager, error) {
	b, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return nil, err
	}

	return &Manager{
		api:       b,
		usersRepo: usersRepo,
	}, nil
}
