package dto

type RegisterPayload struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,max=666"`
	LastName  string `json:"lastName" validate:"required,max=666"`
	Password  string `json:"password" validate:"required,min=6,max=12"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

type LoginSuccessData struct {
	Token string `json:"token"`
}
