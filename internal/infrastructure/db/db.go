package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// DB ..
type DB struct {
	*sql.DB
}

func New() (*DB, error) {
	db, err := sql.Open("postgres", "user=beslan dbname=buddymap password=beslan sslmode=disable host=postgres_db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
