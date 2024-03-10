package mock

import "github.com/stretchr/testify/mock"

type HelperPasswordMock struct {
	Mock *mock.Mock
}

func NewHelperPasswordMock() *HelperPasswordMock {
	return &HelperPasswordMock{&mock.Mock{}}
}

func (h *HelperPasswordMock) HashPassword(password string) (string, error) {
	args := h.Mock.Called(password)

	value := args.Get(0)
	if value == "" {
		return "", args.Error(1)
	}

	return value.(string), nil
}

func (h *HelperPasswordMock) CheckPasswordHash(password, hash string) bool {
	args := h.Mock.Called(password, hash)
	return args.Get(0).(bool)
}
