package repository

import (
	"cobaMetrics/app/model/entity"
	"context"
	"database/sql"
)

type IAccountRepository interface {
	Add(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.Account, error)
	Update(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error)
	DeleteByEmail(ctx context.Context, tx *sql.Tx, email string) error
	GetAll(ctx context.Context, tx *sql.Tx, limit int, offset int) ([]entity.Account, error)
}
