package types

type NewOptionDTO struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"`
}
