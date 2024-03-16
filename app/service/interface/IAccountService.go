package service

import (
	"cobaMetrics/app/model/dto"
	"context"
)

type IAccountService interface {
	Add(ctx context.Context, request *dto.AddUserRequest) (*dto.AddUserResponse, error)
	GetByEmail(ctx context.Context, email string) (*dto.AccountDetailResponse, error)
	Update(ctx context.Context, request *dto.UpdateAccountRequest) (*dto.AccountDetailResponse, error)
	Login(ctx context.Context, request *dto.LoginRequest) (*dto.LoginResponse, error)
	GetAll(ctx context.Context, limit int, page int) ([]dto.AccountDetailResponse, error)
}
