package customError

import "net/http"

type NotFoundError struct {
	StatusCode int    `json:"status_code"`
	S          string `json:"error"`
}

func NewNotFoundError(s string) error {
	return &NotFoundError{
		StatusCode: http.StatusNotFound,
		S:          s,
	}
}

func (n *NotFoundError) Error() string {
	return n.S
}
