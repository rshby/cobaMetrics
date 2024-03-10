package service

import (
	"cobaMetrics/app/model/dto"
	"context"
)

type IAccountService interface {
	Add(ctx context.Context, request *dto.AddUserRequest) (*dto.AddUserResponse, error)
	GetByEmail(ctx context.Context, email string) (*dto.AccountDetailResponse, error)
}
