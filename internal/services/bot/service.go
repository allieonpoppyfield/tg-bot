package bot

import (
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/cache"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotService struct {
	Cache *cache.Cache
	Repo  repositories.BotRepository
}

func New(api *tgbotapi.BotAPI, cache *cache.Cache, repo repositories.BotRepository) (*BotService, error) {
	return &BotService{
		Cache: cache,
		Repo:  repo,
	}, nil
}
