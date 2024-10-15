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
			StateFllingSurvey: &FillingSurveyState{
				state: StateFllingSurvey,
			},
			StateMain: &MainState{
				state: StateMain,
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
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == string(StateProfile) {
			fsm.steps[StateProfile].Execute(update, fsm.api, fsm.service)
		}
	} else if curr, ok := fsm.steps[State(update.Message.Text)]; ok {
		_, err := fsm.api.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID))
		if err != nil {
			logger.Error(
				"cannot delete message",
				slog.String("message", update.Message.Text),
			)
		}
		curr.Execute(update, fsm.api, fsm.service)
	} else if fsm.service.Cache.Client.Get(ctx, string(chatID)) != nil {
		if slices.Contains([]SessionState{sName, sAge, sGender, sDescription, sPhoto}, SessionState(fsm.service.Cache.Client.Get(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatID)).Val())) {
			fsm.steps[StateFllingSurvey].Execute(update, fsm.api, fsm.service)
		} else {
			fsm.steps[StateMain].Execute(update, fsm.api, fsm.service)
		}
	}
}
