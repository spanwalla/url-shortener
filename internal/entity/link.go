package entity

type Link struct {
	Alias string `db:"alias"`
	URI   string `db:"uri"`
}
