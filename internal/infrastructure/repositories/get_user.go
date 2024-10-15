package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/allieonpoppyfield/tg-bot/internal/infrastructure/models"
)

func (r *botRepository) GetUser(ctx context.Context, chatID int64) (*models.DbUser, error) {

	selectQuery := queryBuilder.Select(
		"name",
		"age",
		"gender",
		"description",
		"photo_url",
	).From("users").Where(squirrel.Eq{"chat_id": chatID})

	query, args, err := selectQuery.ToSql()
	if err != nil {
		return nil, err
	}

	// Объявление переменной для хранения результата
	var user models.DbUser

	// Выполнение запроса и заполнение структуры User
	err = r.db.QueryRow(query, args...).Scan(
		&user.Name,
		&user.Age,
		&user.Gender,
		&user.Description,
		&user.PhotoURL,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
