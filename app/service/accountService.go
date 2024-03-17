package service

import (
	"cobaMetrics/app/config"
	"cobaMetrics/app/customError"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	"cobaMetrics/app/model/entity"
	jwtModel "cobaMetrics/app/model/jwt"
	IRepo "cobaMetrics/app/repository/interface"
	IService "cobaMetrics/app/service/interface"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type AccountService struct {
	DB             *sql.DB
	Validate       *validator.Validate
	AccRepo        IRepo.IAccountRepository
	HelperPassword helper.IHelperPassword
	Config         config.IConfig
}

func NewAccountService(db *sql.DB, validate *validator.Validate, config config.IConfig, accRepo IRepo.IAccountRepository, helperPassword helper.IHelperPassword) IService.IAccountService {
	return &AccountService{
		DB:             db,
		Validate:       validate,
		Config:         config,
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
		span.LogFields(log.String("response", err.Error()))
		ext.Error.Set(span, true)
		return nil, err
	}

	// hash password
	hashedPassword, err := a.HelperPassword.HashPassword(request.Password)
	if err != nil {
		span.LogFields(log.String("response", err.Error()))
		ext.Error.Set(span, true)
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

	// cek if email already exist
	if _, err = a.AccRepo.GetByEmail(ctxTracing, tx, request.Email); err == nil {
		errMessage := "email already exist in database"
		span.LogFields(log.String("response", errMessage))
		ext.Error.Set(span, true)
		return nil, customError.NewBadRequestError(errMessage)
	}

	// call procedure insert in repository
	account, err := a.AccRepo.Add(ctxTracing, tx, &input)
	if err != nil {
		span.LogFields(log.String("response", err.Error()))
		ext.Error.Set(span, true)
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
		span.LogFields(log.String("response", err.Error()))
		ext.Error.Set(span, true)
		return nil, err
	}

	// start transaction
	tx, _ := a.DB.Begin()
	defer tx.Rollback()

	// call procedure in repository
	account, err := a.AccRepo.GetByEmail(ctxTracing, tx, email)
	if err != nil {
		span.LogFields(log.String("response", err.Error()))
		ext.Error.Set(span, true)
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

// implementasi method Update data account
func (a *AccountService) Update(ctx context.Context, request *dto.UpdateAccountRequest) (*dto.AccountDetailResponse, error) {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService Update")
	defer span.Finish()

	// log with tracing
	requestJson, _ := json.Marshal(&request)
	span.LogFields(
		log.String("request", string(requestJson)))

	// validate
	if err := a.Validate.Struct(*request); err != nil {
		ext.Error.Set(span, true)
		return nil, err
	}

	// hashed new password
	hashedPassword, err := a.HelperPassword.HashPassword(request.Password)
	if err != nil {
		ext.Error.Set(span, true)
		return nil, customError.NewInternalServerError(err.Error())
	}

	// create input
	input := entity.Account{
		Id:       request.Id,
		Email:    request.Email,
		Username: request.Username,
		Password: hashedPassword,
	}

	// create transaction
	tx, _ := a.DB.Begin()
	defer tx.Rollback()

	// call procedure in repository
	account, err := a.AccRepo.Update(ctxTracing, tx, &input)
	if err != nil {
		ext.Error.Set(span, true)
		return nil, err
	}

	// create response
	response := dto.AccountDetailResponse{
		Id:        account.Id,
		Email:     account.Email,
		Username:  account.Username,
		Password:  account.Password,
		UpdatedAt: helper.DateToString(account.UpdatedAt),
	}

	// log witn tracing
	responseJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(responseJson)))

	tx.Commit()
	return &response, nil
}

// implementasi method Login
func (a *AccountService) Login(ctx context.Context, request *dto.LoginRequest) (*dto.LoginResponse, error) {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService Login")
	defer span.Finish()

	// logging with tracing
	requestJson, _ := json.Marshal(&request)
	span.LogFields(log.String("request", string(requestJson)))

	// validate
	if err := a.Validate.Struct(*request); err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("response", err.Error()))
		return nil, err
	}

	// start transaction
	tx, _ := a.DB.Begin()
	defer tx.Rollback()

	// cek email
	account, err := a.AccRepo.GetByEmail(ctxTracing, tx, request.Email)
	if err != nil {
		errMessage := "record not found"
		ext.Error.Set(span, true)
		span.LogFields(log.String("response", errMessage))
		return nil, customError.NewNotFoundError(errMessage)
	}

	// check password
	if isValid := a.HelperPassword.CheckPasswordHash(request.Password, account.Password); !isValid {
		errMessage := "password not match"
		ext.Error.Set(span, true)
		span.LogFields(log.String("response", errMessage))
		return nil, customError.NewBadRequestError(errMessage)
	}

	// create token
	jwtConfig := a.Config.Config().Jwt
	claims := jwtModel.Claims{
		Email: account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtConfig.Issuer,
			Subject:   jwtConfig.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	token, err := tokenWithClaims.SignedString([]byte(jwtConfig.SecretKey))
	if err != nil {
		span.LogFields(log.String("response", err.Error()))
		ext.Error.Set(span, true)
		return nil, customError.NewInternalServerError(err.Error())
	}

	// success create Token
	// create response
	response := dto.LoginResponse{
		Token:     token,
		CreatedAt: helper.DateToString(time.Now()),
	}

	// log with tracing
	responseJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(responseJson)))

	tx.Commit()
	return &response, nil
}

// implementasi method GetAll
func (a *AccountService) GetAll(ctx context.Context, limit int, page int) ([]dto.AccountDetailResponse, error) {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService GetAll")
	defer span.Finish()

	// log request
	req := map[string]int{
		"limit": limit,
		"page":  page,
	}
	reqJson, _ := json.Marshal(&req)
	span.LogFields(log.String("request", string(reqJson)))

	// create transaction
	tx, _ := a.DB.BeginTx(ctxTracing, nil)
	defer tx.Rollback()

	// call procedure GetAll in repository
	offset := (limit * page) - limit
	accounts, err := a.AccRepo.GetAll(ctxTracing, tx, limit, offset)
	if err != nil {
		span.LogFields(log.String("response", err.Error()))
		return nil, err
	}

	var response []dto.AccountDetailResponse

	wg := &sync.WaitGroup{}
	for _, account := range accounts {
		wg.Add(1)
		go func(wg *sync.WaitGroup, acc entity.Account) {
			defer wg.Done()
			response = append(response, dto.AccountDetailResponse{
				Id:        acc.Id,
				Email:     acc.Email,
				Username:  acc.Username,
				Password:  acc.Password,
				CreatedAt: helper.DateToString(acc.CreatedAt),
				UpdatedAt: helper.DateToString(acc.UpdatedAt),
			})
		}(wg, account)
	}

	wg.Wait()

	// log response
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))

	tx.Commit()
	return response, nil
}
