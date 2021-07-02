package types

import "github.com/jackc/pgtype"

type NewUserDTO struct {
	ID        int                `db:"id"`
	Name      string             `db:"name"`
	Email     string             `db:"email"`
	Password  string             `db:"password"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	DeletedAt pgtype.Timestamptz `db:"deleted_at"`
}

type NewUserForm struct {
	ID                   int
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
}
