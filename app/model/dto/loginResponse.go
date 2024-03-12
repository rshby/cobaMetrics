package dto

type LoginResponse struct {
	Token     string `json:"token,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
