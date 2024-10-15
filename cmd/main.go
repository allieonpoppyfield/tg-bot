package main

import (
	"context"
	"log"

	"github.com/allieonpoppyfield/tg-bot/config"
	"github.com/allieonpoppyfield/tg-bot/internal/fsm"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/cache"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/db"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/repositories"
	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	"github.com/allieonpoppyfield/tg-bot/internal/services/yandex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	ctx := context.TODO()
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	api, err := tgbotapi.NewBotAPI(config.Bot.Key)
	if err != nil {
		log.Fatal(err)
	}
	cache, err := cache.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	repo := repositories.New(db)
	botService, err := bot.New(api, cache, repo, yandex.New())
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	manager := fsm.NewFSMManager(botService, api)
	if err = manager.Start(); err != nil {
		log.Fatal(err)
	}
}
