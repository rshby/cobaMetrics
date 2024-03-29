package repository

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/model/entity"
	IRepo "cobaMetrics/app/repository/interface"
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type AccountRepository struct {
}

// function provider
func NewAccountRepository() IRepo.IAccountRepository {
	return &AccountRepository{}
}

// method implementasi Add new data account
func (a *AccountRepository) Add(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error) {
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
	input.CreatedAt = time.Now()

	// success
	responseJson, _ := json.Marshal(&input)
	span.LogFields(log.String("response", string(responseJson)))
	return input, nil
}

// method implementasi GetByEmail
func (a *AccountRepository) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.Account, error) {
	// span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountRepository Get By Email")
	defer span.Finish()

	// log with tracer
	span.LogFields(
		log.String("email", email))

	// execute
	row := tx.QueryRowContext(ctxTracing, "SELECT id, email, username, password, created_at, updated_at FROM accounts WHERE email = ?", email)
	if row.Err() != nil {
		return nil, customError.NewInternalServerError(row.Err().Error())
	}

	// scan
	account := entity.Account{}
	if err := row.Scan(&account.Id, &account.Email, &account.Username, &account.Password, &account.CreatedAt, &account.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, customError.NewNotFoundError("record not found")
		}

		return nil, customError.NewInternalServerError(err.Error())
	}

	// success
	// log with tracer
	responseJson, _ := json.Marshal(&account)
	span.LogFields(
		log.String("response", string(responseJson)))

	return &account, nil
}

// implementasi method Update data account
func (a *AccountRepository) Update(ctx context.Context, tx *sql.Tx, input *entity.Account) (*entity.Account, error) {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountRepository Update")
	defer span.Finish()

	// log with tracing
	requestJson, _ := json.Marshal(&input)
	span.LogFields(
		log.String("request", string(requestJson)))

	// update
	result, err := tx.ExecContext(ctxTracing, "UPDATE accounts SET email=?, username=?, password=? WHERE id = ?",
		input.Email, input.Username, input.Password, input.Id)
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return nil, customError.NewInternalServerError("failed to update data account")
	}

	input.UpdatedAt = time.Now()

	// log with tracing
	responseJson, _ := json.Marshal(&input)
	span.LogFields(
		log.String("response", string(responseJson)))

	return input, nil
}

func (a *AccountRepository) DeleteByEmail(ctx context.Context, tx *sql.Tx, email string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepository) GetAll(ctx context.Context, tx *sql.Tx, limit int, offset int) ([]entity.Account, error) {
	// start tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountRepository GetAll")
	defer span.Finish()

	// log with tracing
	req := map[string]int{
		"limit":  limit,
		"offset": offset,
	}

	reqJson, _ := json.Marshal(&req)
	span.LogFields(log.String("request", string(reqJson)))

	// query
	rows, err := tx.QueryContext(ctxTracing, "SELECT id, email, username, password, created_at, updated_at FROM accounts ORDER BY accounts.id LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		// log error
		span.LogFields(log.String("response", err.Error()))
		return nil, customError.NewInternalServerError(err.Error())
	}

	if rows.Err() != nil {
		// log with tracing
		span.LogFields(log.String("response", rows.Err().Error()))
		return nil, customError.NewInternalServerError(rows.Err().Error())
	}

	var accounts []entity.Account
	for rows.Next() {
		var account entity.Account
		if err = rows.Scan(&account.Id, &account.Email, &account.Username, &account.Password, &account.CreatedAt, &account.UpdatedAt); err != nil {
			span.LogFields(log.String("response", err.Error()))

			// if error not found data
			if err == sql.ErrNoRows {
				return nil, customError.NewNotFoundError(err.Error())
			}

			return nil, customError.NewInternalServerError(err.Error())
		}

		// append to accounts
		accounts = append(accounts, account)
	}

	// log response with tracing
	resJson, _ := json.Marshal(&accounts)
	span.LogFields(log.String("response", string(resJson)))

	return accounts, nil
}
