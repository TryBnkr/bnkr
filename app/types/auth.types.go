package types

import "database/sql"

// LoginDTO defined the /login payload
type LoginDTO struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=6"`
}

// UserResponse
type UserResponse struct {
	ID        int          `json:"id" db:"id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Name      string       `json:"name" db:"name"`
	Email     string       `json:"email" db:"email"`
	Password  string       `json:"-" db:"password"`
}
