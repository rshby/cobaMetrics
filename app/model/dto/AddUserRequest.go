package dto

type AddUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=6"`
}
