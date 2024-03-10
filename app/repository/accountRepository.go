package repository

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/model/entity"
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Account struct {
	DB *sql.DB
}

func (a *Account) Add(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error) {
	// tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Repository Add Account")
	defer span.Finish()

	inputJson, _ := json.Marshal(&input)
	span.LogFields(log.String("request", string(inputJson)))

	query := "INSERT INTO accounts(email, username, password) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctxTracing, query, input.Email, input.Username, input.Password)
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return nil, customError.NewInternalServerError("failed to insert new user")
	}

	id, _ := result.LastInsertId()
	input.Id = int(id)

	// success
	responseJson, _ := json.Marshal(&input)
	span.LogFields(log.String("response", string(responseJson)))
	return input, nil
}

func (a *Account) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Account) Update(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Account) DeleteByEmail(ctx context.Context, tx *sql.Tx, email string) error {
	//TODO implement me
	panic("implement me")
}

func (a *Account) GetAll(ctx context.Context, tx *sql.Tx, limit int, offset int) ([]entity.Account, error) {
	//TODO implement me
	panic("implement me")
}
