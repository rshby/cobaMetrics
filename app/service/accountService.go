package service

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	"cobaMetrics/app/model/entity"
	IRepo "cobaMetrics/app/repository/interface"
	IService "cobaMetrics/app/service/interface"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type AccountService struct {
	DB             *sql.DB
	Validate       *validator.Validate
	AccRepo        IRepo.IAccountRepository
	HelperPassword helper.IHelperPassword
}

func NewAccountService(db *sql.DB, validate *validator.Validate, accRepo IRepo.IAccountRepository, helperPassword helper.IHelperPassword) IService.IAccountService {
	return &AccountService{
		DB:             db,
		Validate:       validate,
		AccRepo:        accRepo,
		HelperPassword: helperPassword,
	}
}

// method implementasi Add new account
func (a *AccountService) Add(ctx context.Context, request *dto.AddUserRequest) (*dto.AddUserResponse, error) {
	// start tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Service Add Account")
	defer span.Finish()

	// log request
	reqJson, _ := json.Marshal(&request)
	span.LogFields(log.String("request", string(reqJson)))

	// validate
	if err := a.Validate.Struct(*request); err != nil {
		return nil, err
	}

	// hash password
	hashedPassword, err := a.HelperPassword.HashPassword(request.Password)
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	// create input
	input := entity.Account{
		Email:     request.Email,
		Username:  request.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// begin transaction
	tx, _ := a.DB.Begin()
	defer tx.Rollback()

	// call procedure insert in repository
	account, err := a.AccRepo.Add(ctxTracing, tx, &input)
	if err != nil {
		return nil, err
	}

	// create response
	response := dto.AddUserResponse{
		Id:        account.Id,
		Email:     account.Email,
		Username:  account.Username,
		CreatedAt: account.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// success
	tx.Commit()
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))
	return &response, nil
}

// method implementasi GetByEmail
func (a *AccountService) GetByEmail(ctx context.Context, email string) (*dto.AccountDetailResponse, error) {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService GetByEmail")
	defer span.Finish()

	// log with tracer
	span.LogFields(
		log.String("email", email))

	// validate
	if err := a.Validate.Var(email, "email"); err != nil {
		return nil, err
	}

	// start transaction
	tx, _ := a.DB.Begin()
	defer tx.Rollback()

	// call procedure in repository
	account, err := a.AccRepo.GetByEmail(ctxTracing, tx, email)
	if err != nil {
		return nil, err
	}

	// convert to response
	var response = dto.AccountDetailResponse{
		Id:        account.Id,
		Email:     account.Email,
		Username:  account.Username,
		Password:  account.Password,
		CreatedAt: helper.DateToString(account.CreatedAt),
		UpdatedAt: helper.DateToString(account.UpdatedAt),
	}

	// log with tracer
	responseJson, _ := json.Marshal(&response)
	span.LogFields(
		log.String("response", string(responseJson)))

	tx.Commit()
	return &response, nil
}
