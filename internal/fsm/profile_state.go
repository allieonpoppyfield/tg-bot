package fsm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ProfileState struct {
	state              State
	previousState      State
	availableStateList []State
}

type User struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Gender      int16  `json:"gender"`
	Description string `json:"description"`
}

func (s *ProfileState) Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	ctx := context.Background()
	chatId := update.Message.Chat.ID

	// Получаем текущее состояние из кеша
	cacheVal, err := service.Cache.Client.Get(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId)).Result()
	if err != nil {
		// Если нет состояния, начинаем с имени
		return s.askName(ctx, chatId, api, service)
	}

	// Обрабатываем состояние на основе текущего значения
	return s.handleState(ctx, chatId, update.Message.Text, cacheVal, api, service)
}

func (s *ProfileState) askName(ctx context.Context, chatId int64, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	// Сохраняем текущее состояние
	service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sName), 5*time.Minute)
	msg := tgbotapi.NewMessage(chatId, "Как вас зовут?")
	api.Send(msg)
	return StateProfile, nil
}

func (s *ProfileState) handleState(ctx context.Context, chatId int64, userInput string, stateCache string, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	var u User
	uStr := service.Cache.Client.Get(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId)).Val()
	json.Unmarshal([]byte(uStr), &u) // Десериализуем пользователя

	switch SessionState(stateCache) {
	case sName:
		u.Name = userInput
		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sAge), 5*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Сколько вам лет?"))
	case sAge:
		age, err := strconv.Atoi(userInput)
		if err != nil || age < 14 || age > 99 {
			api.Send(tgbotapi.NewMessage(chatId, "Введите корректный возраст"))
			return StateProfile, nil
		}
		u.Age = age
		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sGender), 5*time.Minute)
		service.Cache.Client.Set(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId), marshalUser(u), 5*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Какого вы пола? (м / ж)"))
	case sGender:
		if !strings.EqualFold(userInput, "м") && !strings.EqualFold(userInput, "ж") {
			api.Send(tgbotapi.NewMessage(chatId, "Нужно указать ваш пол (м / ж)"))
			return StateProfile, nil
		}
		u.Gender = 1
		if strings.EqualFold(userInput, "ж") {
			u.Gender = 2
		}
		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sDescription), 5*time.Minute)
		service.Cache.Client.Set(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId), marshalUser(u), 5*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Расскажите немного о себе"))
	case sDescription:
		u.Description = userInput
		if err := service.Repo.InsertUser(ctx, chatId, u.Name, u.Age, u.Gender, u.Description); err != nil {
			api.Send(tgbotapi.NewMessage(chatId, err.Error()))
			return StateMain, nil
		}
		// Удаляем кэшированные данные
		service.Cache.Client.Del(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId))
		service.Cache.Client.Del(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId))
		api.Send(tgbotapi.NewMessage(chatId, "Профиль успешно сохранен!"))
	}

	return StateMain, nil
}

// Упрощаем маршалинг объекта User
func marshalUser(u User) []byte {
	um, _ := json.Marshal(u)
	return um
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
