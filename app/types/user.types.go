package types

type NewUserDTO struct {
	ID                   uint
	Name                 string
	Email                string
	Password             string
}

type NewUserForm struct {
	ID                   uint
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
}
