package fsm

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"slices"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"

	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type FSMManager struct {
	service *bot.BotService
	api     *tgbotapi.BotAPI
	steps   map[State]FSMState
}

func NewFSMManager(s *bot.BotService, a *tgbotapi.BotAPI) *FSMManager {
	return &FSMManager{
		service: s,
		api:     a,
		steps: map[State]FSMState{
			StateStart: &StartState{
				state:              StateStart,
				availableStateList: workflow[StateStart],
			},
			StateProfile: &ProfileState{
				state: StateProfile,
			},
		},
	}
}

func (fsm *FSMManager) Start() error {
	zapL := zap.Must(zap.NewProduction())
	defer zapL.Sync()

	fsm.api.Debug = true
	log.Printf("Authorized on account %s", fsm.api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := fsm.api.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}
	// Обработка обновлений
	for update := range updates {
		fsm.handle(update)
	}
	return nil
}

func (fsm *FSMManager) handle(update tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	zapL := zap.Must(zap.NewProduction())
	defer zapL.Sync()
	logger := slog.New(zapslog.NewHandler(zapL.Core(), nil))
	chatID := update.Message.Chat.ID
	if curr, ok := fsm.steps[State(update.Message.Text)]; ok {
		_, err := fsm.api.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))
		if err != nil {
			logger.Error(
				"cannot delete message",
				slog.String("message", update.Message.Text),
			)
		}
		curr.Execute(update, fsm.api, fsm.service)
	} else if fsm.service.Cache.Client.Get(ctx, string(chatID)) != nil {
		if slices.Contains([]SessionState{sName, sAge, sGender, sDescription}, SessionState(fsm.service.Cache.Client.Get(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatID)).Val())) {
			fsm.steps[StateProfile].Execute(update, fsm.api, fsm.service)
		}
	} else {
		logger.Error(
			"next step not found",
			slog.String("step", update.Message.Text),
		)
	}

	// // Инициализация состояния, если не существует
	// if _, exists := userStates[chatID]; !exists {
	// 	userStates[chatID] = &User{
	// 		State: State(), // Начальное состояние
	// 	}
	// 	bot.SendMessage(chatID, "Введите свое имя:")
	// 	return
	// }

	// // Обработка текущего состояния пользователя
	// user := userStates[chatID]
	// step := fsm.steps[user.State]

	// if step != nil {
	// 	nextState, err := step.Execute(update)
	// 	if err != nil {
	// 		log.Printf("Error executing state %s: %v", user.State, err)
	// 	}

	// 	// Переход на следующее состояние, полученное из Execute
	// 	user.State = nextState
	// } else {
	// 	bot.SendMessage(chatID, "Неизвестное состояние.")
	// }
}
