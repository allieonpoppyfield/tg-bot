package repositories

import (
	"context"
	"time"
)

func (r *botRepository) InsertUser(ctx context.Context, chatID int64, name string, age int, gender int16, desc string) error {
	params := map[string]any{
		"name":        name,
		"chat_id":     chatID,
		"age":         age,
		"gender":      gender,
		"description": desc,
		"created_at":  time.Now().UTC(),
	}

	insertQuery := queryBuilder.Insert("users").SetMap(params)

	query, args, err := insertQuery.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}
