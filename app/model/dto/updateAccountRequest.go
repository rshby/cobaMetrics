package dto

type UpdateAccountRequest struct {
	Id              int    `json:"id,omitempty" validate:"required,gt=0"`
	Email           string `json:"email,omitempty" validate:"required,email"`
	Username        string `json:"username,omitempty" validate:"required,min=2"`
	Password        string `json:"password,omitempty" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password,omitempty" validate:"required,eqfield=Password"`
}
