package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BroadcastWbCatalogNotificationParam struct {
	ImgURL, NotificationName, ProductURL string
	Price, Quantity                      int64
}

func (m *Manager) BroadcastWbCatalogNotification(params BroadcastWbCatalogNotificationParam) error {
	users, err := m.usersRepo.GetAll()
	if err != nil {
		return err
	}

	for _, user := range users {
		msg := tgbotapi.NewPhoto(user.UserID, tgbotapi.FileURL(params.ImgURL))
		caption := "<strong>WB новый товар в каталоге</strong>\n\n"
		caption += fmt.Sprintf("<strong>Категория</strong>: %s\n", params.NotificationName)
		caption += fmt.Sprintf("<strong>Цена</strong>: %d\n", params.Price)
		caption += fmt.Sprintf("<strong>Количество</strong>: %d", params.Quantity)
		msg.Caption = caption

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Открыть на сайте", params.ProductURL),
			),
		)

		msg.ReplyMarkup = inlineKeyboard
		msg.ParseMode = "HTML"

		_, err := m.api.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
