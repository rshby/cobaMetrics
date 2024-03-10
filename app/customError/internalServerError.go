package customError

import "net/http"

type InternalServerError struct {
	StatusCode int    `json:"status_code"`
	S          string `json:"error"`
}

func NewInternalServerError(s string) error {
	return &InternalServerError{
		StatusCode: http.StatusInternalServerError,
		S:          s,
	}
}

func (i *InternalServerError) Error() string {
	return i.S
}
