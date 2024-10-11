package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Хендлер для команды /location
func LocationHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	// Запрос на отправку локации
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, поделитесь своей локацией.")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("Отправить локацию"),
		),
	)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// Обработка полученной локации
func HandleLocation(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Location != nil {
		lat := update.Message.Location.Latitude
		lon := update.Message.Location.Longitude

		log.Printf("Получена локация: %.6f, %.6f", lat, lon)

		// Здесь можно реализовать логику обновления локации в БД или отправки на сервер
		// Например, можно использовать сервис для сохранения локации
		// services.SaveLocation(update.Message.From.ID, lat, lon)

		// Отправка подтверждения
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Локация получена!")
		bot.Send(msg)
	}
}
