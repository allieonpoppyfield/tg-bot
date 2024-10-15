package fsm

import (
	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// func (s *MainState) Execute(update tgbotapi.Update, bot *bot.BotService) (State, error) {
// 	chatID := update.Message.Chat.ID
// 	// Создаем кнопки
// 	keyboard := tgbotapi.NewReplyKeyboard(
// 		tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("Оставить заявку"),
// 		),
// 	)

// 	// Отправляем сообщение с клавиатурой
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
// 	msg.ReplyMarkup = keyboard

// 	bot.ap.SendMessage(msg)
// 	return StateAge, nil
// }

type MainState struct {
	state              State
	previousState      State
	availableStateList []State
}

func (s *MainState) Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	// Создание инлайн-кнопок
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мой профиль", "profile"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	api.Send(msg)
	return StateMain, nil
}

func (s *MainState) GetState() State {
	return s.state
}

func (s *MainState) SetPreviousState(st State) {
	s.previousState = st
}

func (s *MainState) GetPreviousState() State {
	return s.previousState
}

func (s *MainState) GetAvailableStateList() []State {
	return s.availableStateList
}
