package telegrambotusers

import (
	"context"
	"core-consumer/internal/app/gen/model"
)

func (r *Repository) GetAll() ([]*model.TelegramBotUser, error) {
	return r.q.TelegramBotUser.WithContext(context.TODO()).Find()
}
