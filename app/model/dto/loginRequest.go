package dto

type LoginRequest struct {
	Email    string `json:"username,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}
