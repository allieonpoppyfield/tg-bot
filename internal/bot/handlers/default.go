package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// DefaultHandler - обработчик по умолчанию для сообщений
func DefaultHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Проверяем, что это сообщение, а не обновление типа 'inline'
	if update.Message == nil {
		return
	}

	// Логируем текст сообщения
	log.Printf("Received message from %s: %s", update.Message.From.UserName, update.Message.Text)

	// Формируем ответное сообщение
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, я не понимаю эту команду. Попробуйте /help для получения списка доступных команд.")

	// Отправляем ответ
	if _, err := bot.Send(reply); err != nil {
		log.Printf("Error sending message: %s", err)
	}
}
