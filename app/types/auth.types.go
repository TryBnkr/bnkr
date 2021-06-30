package types

import "github.com/jackc/pgtype"

// LoginDTO defined the /login payload
type LoginDTO struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=6"`
}

// UserResponse
type UserResponse struct {
	ID        int                `json:"id" db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	DeletedAt pgtype.Timestamptz `db:"deleted_at"`
	Name      string             `json:"name" db:"name"`
	Email     string             `json:"email" db:"email"`
	Password  string             `json:"-" db:"password"`
}
