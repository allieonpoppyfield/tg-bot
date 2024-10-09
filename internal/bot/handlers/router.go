package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Router struct {
	handlers map[string]func(*tgbotapi.BotAPI, tgbotapi.Update)
}

func NewRouter() *Router {
	r := &Router{
		handlers: make(map[string]func(*tgbotapi.BotAPI, tgbotapi.Update)),
	}
	r.registerHandlers()
	return r
}

func (r *Router) registerHandlers() {
	r.handlers["/start"] = StartHandler
	r.handlers["/location"] = LocationHandler // Добавляем хендлер для /location
}

func (r *Router) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.Location != nil {
			// Обработка полученной локации
			HandleLocation(bot, update)
			return
		}

		// Если это текстовое сообщение, то ищем соответствующий хендлер
		handler, exists := r.handlers[update.Message.Text]
		if exists {
			handler(bot, update)
		} else {
			// Хендлер по умолчанию
			DefaultHandler(bot, update)
		}
	}
}
