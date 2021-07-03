package types

import "database/sql"

type NewUserDTO struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type NewUserForm struct {
	ID                   int
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
}
