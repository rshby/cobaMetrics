package mock

import (
	"cobaMetrics/app/model/entity"
	"context"
	"database/sql"
	"github.com/stretchr/testify/mock"
)

type AccountRepositoryMock struct {
	Mock *mock.Mock
}

func NewAccountRepository() *AccountRepositoryMock {
	return &AccountRepositoryMock{&mock.Mock{}}
}

func (a *AccountRepositoryMock) Add(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error) {
	args := a.Mock.Called(ctx, tx, input)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	// not error
	return value.(*entity.Account), nil
}

func (a *AccountRepositoryMock) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepositoryMock) Update(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepositoryMock) DeleteByEmail(ctx context.Context, tx *sql.Tx, email string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepositoryMock) GetAll(ctx context.Context, tx *sql.Tx, limit int, offset int) ([]entity.Account, error) {
	//TODO implement me
	panic("implement me")
}
