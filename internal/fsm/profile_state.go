package fsm

import (
	"context"
	"fmt"
	"log"

	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ProfileState struct {
	state              State
	previousState      State
	availableStateList []State
}

func (s *ProfileState) Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	user, err := service.Repo.GetUser(context.Background(), update.Message.Chat.ID)
	if err != nil {
		return StateMain, err
	}
	text := fmt.Sprintf("Имя: %s\nВозраст: %d\nПол: %s\nОписание: %s", user.Name, user.Age, user.Gender, user.Description)
	err = service.Ya.DownloadPhoto(user.PhotoURL, string(update.Message.Chat.ID)+".jpg")
	if err != nil {
		log.Fatal(err)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	api.Send(msg)

	return StateMain, nil
}

func (s *ProfileState) GetState() State {
	return s.state
}

func (s *ProfileState) SetPreviousState(st State) {
	s.previousState = st
}

func (s *ProfileState) GetPreviousState() State {
	return s.previousState
}

func (s *ProfileState) GetAvailableStateList() []State {
	return s.availableStateList
}
