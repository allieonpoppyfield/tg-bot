package fsm

import (
	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type State string

const (
	StateStart   State = "/start"
	StateProfile State = "Заполнить анкету"
	StateMain    State = "main"
)

type SessionState string

const (
	sName        SessionState = "name"
	sAge         SessionState = "age"
	sGender      SessionState = "gender"
	sDescription SessionState = "description"
)

const SESSION_STATE_CACHE_KEY = "session_state %d:"
const PROFILE_CACHE_KEY = "profile %d:"

type FSMState interface {
	Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) // Возвращаем следующее состояние в зависимости от условий
	GetState() State
	GetPreviousState() State
	SetPreviousState(State)
	GetAvailableStateList() []State
}

var workflow map[State][]State = map[State][]State{
	StateStart: {
		StateMain,
		StateProfile,
	},
}
