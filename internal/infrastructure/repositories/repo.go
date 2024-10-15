package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/db"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/models"
)

var queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type BotRepository interface {
	InsertUser(ctx context.Context, chatID int64, name string, age int, gender int16, desc string, photo_url string) error
	GetUser(ctx context.Context, chatID int64) (*models.DbUser, error)
}

type botRepository struct {
	db *db.DB
}

// New ..
func New(db *db.DB) BotRepository {
	return &botRepository{
		db: db,
	}
}
