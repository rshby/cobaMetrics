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

func (a *AccountServiceMock) GetByEmail(ctx context.Context, email string) (*dto.AccountDetailResponse, error) {
	args := a.Mock.Called(ctx, email)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*dto.AccountDetailResponse), nil
}

func (a *AccountServiceMock) Update(ctx context.Context, request *dto.UpdateAccountRequest) (*dto.AccountDetailResponse, error) {
	args := a.Mock.Called(ctx, request)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*dto.AccountDetailResponse), nil
}

func (a *AccountServiceMock) Login(ctx context.Context, request *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := a.Mock.Called(ctx, request)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*dto.LoginResponse), nil
}

func (a *AccountServiceMock) GetAll(ctx context.Context, limit int, page int) ([]dto.AccountDetailResponse, error) {
	args := a.Mock.Called(ctx, limit, page)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.([]dto.AccountDetailResponse), nil
}
