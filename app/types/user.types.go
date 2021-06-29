package types

type NewUserDTO struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type NewUserForm struct {
	ID                   int
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
}
