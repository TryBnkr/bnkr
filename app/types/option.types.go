package types

import "database/sql"

type NewOptionDTO struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Value     string       `db:"value"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
