package models

type DbUser struct {
	Name        string `db:"name"`
	Age         int    `db:"age"`
	Gender      string `db:"gender"`
	Description string `db:"description"`
	PhotoURL    string `db:"photo_url"`
}
