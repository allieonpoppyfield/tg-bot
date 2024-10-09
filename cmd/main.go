package main

import (
	"log"

	"github.com/allieonpoppyfield/tg-bot/config"
	"github.com/allieonpoppyfield/tg-bot/internal/bot"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = bot.Start(config)
	if err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
}
