package customError

import "net/http"

type BadRequestError struct {
	StatusCode int    `json:"status_code"`
	S          string `json:"error"`
}

func NewBadRequestError(s string) error {
	return &BadRequestError{
		StatusCode: http.StatusBadRequest,
		S:          s,
	}
}

func (b *BadRequestError) Error() string {
	return b.S
}
