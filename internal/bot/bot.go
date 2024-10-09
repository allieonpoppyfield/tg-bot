package bot

import (
	"log"

	"github.com/allieonpoppyfield/tg-bot/config"
	"github.com/allieonpoppyfield/tg-bot/internal/bot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Key)
	if err != nil {
		return err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Получение обновлений
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)

	// Регистрация хендлеров
	router := handlers.NewRouter()

	// Обработка обновлений
	for update := range updates {
		router.HandleUpdate(bot, update)
	}

	return nil
}
