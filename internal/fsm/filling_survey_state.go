package fsm

import (
	"context"
	"fmt"
	"log"

	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type FillingSurveyState struct {
	state              State
	previousState      State
	availableStateList []State
}

func (s *FillingSurveyState) Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	ctx := context.Background()
	chatId := update.Message.Chat.ID

	// Получаем текущее состояние из кеша
	cacheVal, err := service.Cache.Client.Get(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId)).Result()
	if err != nil {
		return s.askName(ctx, chatId, api, service)
	}

	// Обработка состояния анкеты
	switch cacheVal {
	case string(sName):
		return s.askAge(ctx, chatId, api, service)
	case string(sAge):
		return s.askGender(ctx, chatId, api, service)
	// Другие шаги анкеты...
	default:
		return s.askName(ctx, chatId, api, service)
	}
}

func (s *FillingSurveyState) askName(ctx context.Context, chatId int64, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	msg := tgbotapi.NewMessage(chatId, "Введите ваше имя:")
	_, err := api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return StateStart, err
	}

	// Сохраняем состояние в кэше
	if err := service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), sName, 0).Err(); err != nil {
		return StateStart, err
	}

	return StateFillingSurvey, nil
}

func (s *FillingSurveyState) askAge(ctx context.Context, chatId int64, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	msg := tgbotapi.NewMessage(chatId, "Введите ваш возраст:")
	_, err := api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return StateStart, err
	}

	// Обновляем состояние в кэше
	if err := service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), sAge, 0).Err(); err != nil {
		return StateStart, err
	}

	return StateFillingSurvey, nil
}

// Дополнительные функции для остальных шагов анкеты (gender, photo, description)
