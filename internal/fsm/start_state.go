package fsm

import (
	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type StartState struct {
	state              State
	previousState      State
	availableStateList []State
}

func (s *StartState) Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствие тырыпыры")
	markup := tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton(string(StateFllingSurvey)),
			},
		},
	}
	msg.ReplyMarkup = markup
	api.Send(msg)
	return StateMain, nil
}

func (s *StartState) GetState() State {
	return s.state
}

func (s *StartState) SetPreviousState(st State) {
	s.previousState = st
}

func (s *StartState) GetPreviousState() State {
	return s.previousState
}

func (s *StartState) GetAvailableStateList() []State {
	return s.availableStateList
}
