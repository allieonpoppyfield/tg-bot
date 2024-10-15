package bot

import (
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/cache"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/repositories"
	"github.com/allieonpoppyfield/tg-bot/internal/services/yandex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotService struct {
	Cache *cache.Cache
	Repo  repositories.BotRepository
	Ya    *yandex.YandexService
}

func New(api *tgbotapi.BotAPI, cache *cache.Cache, repo repositories.BotRepository, ya *yandex.YandexService) (*BotService, error) {
	return &BotService{
		Cache: cache,
		Repo:  repo,
		Ya:    ya,
	}, nil
}
