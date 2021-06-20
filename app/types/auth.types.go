package types

// LoginDTO defined the /login payload
type LoginDTO struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=6"`
}

// JsonLoginDTO defined the /login payload
type JsonLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"password"`
}

// SignupDTO defined the /login payload
type SignupDTO struct {
	JsonLoginDTO
	Name string `json:"name" validate:"required,min=3"`
}

// UserResponse
type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// AccessResponse
type AccessResponse struct {
	Token string `json:"token"`
}

// AuthResponse
type AuthResponse struct {
	User *UserResponse   `json:"user"`
	Auth *AccessResponse `json:"auth"`
}
