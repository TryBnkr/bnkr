package types

import "github.com/jackc/pgtype"

type NewOptionDTO struct {
	ID        int                `db:"id"`
	Name      string             `db:"name"`
	Value     string             `db:"value"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	DeletedAt pgtype.Timestamptz `db:"deleted_at"`
}
