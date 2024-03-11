package mock

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type FieldErrorMock struct {
	validator.FieldError
	TagError string
	FieldErr string
}

func (e *FieldErrorMock) Tag() string { return e.TagError }

func (e *FieldErrorMock) Field() string { return e.FieldErr }

func (e *FieldErrorMock) Error() string {
	return fmt.Sprintf("error on field [%v] with tag [%v]", e.FieldErr, e.TagError)
}
