package mock

import (
	"cobaMetrics/app/model/dto"
	"context"
	"github.com/stretchr/testify/mock"
)

type AccountServiceMock struct {
	Mock *mock.Mock
}

func NewAccountServiceMock() *AccountServiceMock {
	return &AccountServiceMock{&mock.Mock{}}
}

func (a *AccountServiceMock) Add(ctx context.Context, request *dto.AddUserRequest) (*dto.AddUserResponse, error) {
	args := a.Mock.Called(ctx, request)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*dto.AddUserResponse), nil
}
