package test

import (
	"cobaMetrics/app/config"
	"cobaMetrics/app/customError"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	"cobaMetrics/app/model/entity"
	"cobaMetrics/app/service"
	mckConfig "cobaMetrics/app/test/mock/config"
	mckHelper "cobaMetrics/app/test/mock/helper"
	mck "cobaMetrics/app/test/mock/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// unit test method Add
func TestAddUserService(t *testing.T) {
	t.Run("add account error validate", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		validate := validator.New()
		config := mckConfig.NewConfigMock()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountRepository := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, config, accountRepository, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		// test
		request := dto.AddUserRequest{
			Email:    "1",
			Username: "1",
			Password: "12",
		}

		add, err := accountService.Add(context.Background(), &request)
		assert.Nil(t, add)
		assert.NotNil(t, err)
		assert.Error(t, err)

		fieldErrors, ok := err.(validator.ValidationErrors)
		assert.True(t, ok)

		for _, errField := range fieldErrors {
			errMessage := fmt.Sprintf("erro on field [%v] with tag [%v]", errField.Field(), errField.Tag())
			fmt.Println(errMessage)
		}
	})
	t.Run("add account error hash password", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "cant hashing password"
		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("", errors.New(errMessage))

		// test
		request := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		result, err := accountService.Add(context.Background(), &request)
		_, ok := err.(*customError.InternalServerError)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errMessage, err.Error())
		assert.True(t, ok)
		helperPasswordMock.Mock.AssertExpectations(t)
	})
	t.Run("add account error internal server error", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		validate := validator.New()
		helperPassword := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPassword)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPassword.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("record not found")).Times(1)

		errMessage := "failed to add new data"
		accountRepositoryMock.Mock.On("Add", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage))

		// test
		request := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		result, err := accountService.Add(context.Background(), &request)
		_, ok := err.(*customError.InternalServerError)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
		assert.True(t, ok)
		helperPassword.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("add account error bad request", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPassword := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPassword)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPassword.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("record not found")).Times(1)

		errMessage := "error bad request when add data"
		accountRepositoryMock.Mock.On("Add", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewBadRequestError(errMessage))

		// test
		request := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		result, err := accountService.Add(context.Background(), &request)
		_, ok := err.(*customError.BadRequestError)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		helperPassword.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("add account error not found", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("record not found")).Times(1)

		errMessage := "record not found"
		accountRepositoryMock.Mock.On("Add", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage))

		// test
		request := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}

		result, err := accountService.Add(context.Background(), &request)
		_, ok := err.(*customError.NotFoundError)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("add account error email already exist", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil).Times(1)

		errMessage := "email already exist in database"
		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Account{
				Id:        1,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			}, nil).Times(1)

		// test
		account, err := accountService.Add(context.Background(), &dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		})
		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("add account success", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("record not found")).Times(1)

		accountRepositoryMock.Mock.On("Add", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Account{
				Id:        1,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: helper.StringToDate("2020-10-10 00:00:00"),
				UpdatedAt: helper.StringToDate("2020-10-10 00:00:00"),
			}, nil)

		// test
		request := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		result, err := accountService.Add(context.Background(), &request)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "reoshby@gmail.com", result.Email)
		assert.Equal(t, "rshby", result.Username)
		assert.Equal(t, "2020-10-10 00:00:00", result.CreatedAt)
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
}

// unit test method GetByEmail
func TestGetAccountByEmailService(t *testing.T) {
	t.Run("test get account by email error validation", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// test
		email := "reoshby"
		account, err := accountService.GetByEmail(context.Background(), email)
		_, ok := err.(validator.ValidationErrors)
		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
	})
	t.Run("test get account by email error internal server error", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "error internal server error"
		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage))

		// test
		account, err := accountService.GetByEmail(context.Background(), "reoshby@gmail.com")
		_, ok := err.(*customError.InternalServerError)
		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test get account by email error bad request", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "failed to get record error bad request"
		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewBadRequestError(errMessage))

		// test
		account, err := accountService.GetByEmail(context.Background(), "reoshby@gmail.com")
		_, ok := err.(*customError.BadRequestError)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test get account by email error not found", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "record not found"
		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage))

		// test
		account, err := accountService.GetByEmail(context.Background(), "reoshby@gmail.com")
		_, ok := err.(*customError.NotFoundError)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test get account by email success", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Account{
				Id:        1,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: helper.StringToDate("2020-10-10 00:00:00"),
				UpdatedAt: helper.StringToDate("2020-10-10 00:00:00"),
			}, nil)

		// test
		account, err := accountService.GetByEmail(context.Background(), "reoshby@gmail.com")
		assert.Nil(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, 1, account.Id)
		assert.Equal(t, "reoshby@gmail.com", account.Email)
		assert.Equal(t, "rshby", account.Username)
		assert.Equal(t, "2020-10-10 00:00:00", account.CreatedAt)
		assert.Equal(t, "2020-10-10 00:00:00", account.UpdatedAt)
	})
}

// unit test method Update
func TestUpdateAccountService(t *testing.T) {
	t.Run("test update account error validation", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// test
		request := dto.UpdateAccountRequest{
			Id:              0,
			Email:           "reo",
			Username:        "1",
			Password:        "1",
			ConfirmPassword: "3",
		}

		account, err := accountService.Update(context.Background(), &request)
		validationErrors, ok := err.(validator.ValidationErrors)
		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)

		for _, errorField := range validationErrors {
			fmt.Println(fmt.Sprintf("error on field [%v] with tag [%v] : [%v]",
				errorField.Field(), errorField.Tag(), errorField.Error()))
		}
	})
	t.Run("test update account error hash password", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		errMessage := "cant hash password"
		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("", errors.New(errMessage))

		// test
		request := dto.UpdateAccountRequest{
			Id:              1,
			Email:           "reoshby@gmail.com",
			Username:        "rshby",
			Password:        "123456",
			ConfirmPassword: "123456",
		}
		account, err := accountService.Update(context.Background(), &request)
		_, ok := err.(*customError.InternalServerError)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		helperPasswordMock.Mock.AssertExpectations(t)
	})
	t.Run("test update error internal server", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		accountRepositoryMock := mck.NewAccountRepository()
		configMock := mckConfig.NewConfigMock()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		errMessage := "error internal server"
		accountRepositoryMock.Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage))

		// test
		request := dto.UpdateAccountRequest{
			Id:              1,
			Email:           "reoshby@gmail.com",
			Username:        "rshby",
			Password:        "123456",
			ConfirmPassword: "123456",
		}
		account, err := accountService.Update(context.Background(), &request)
		_, ok := err.(*customError.InternalServerError)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test update error not found", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		accountRepositoryMock := mck.NewAccountRepository()
		configMock := mckConfig.NewConfigMock()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		errMessage := "record not fouund"
		accountRepositoryMock.Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage))

		// test
		request := dto.UpdateAccountRequest{
			Id:              1,
			Email:           "reoshby@gmail.com",
			Username:        "rshby",
			Password:        "123456",
			ConfirmPassword: "123456",
		}
		account, err := accountService.Update(context.Background(), &request)
		_, ok := err.(*customError.NotFoundError)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test upate error bad request", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

		errMessage := "error bad request"
		accountRepositoryMock.Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage))

		// test
		request := dto.UpdateAccountRequest{
			Id:              1,
			Email:           "reoshby@gmail.com",
			Username:        "rshby",
			Password:        "123456",
			ConfirmPassword: "123456",
		}
		account, err := accountService.Update(context.Background(), &request)
		_, ok := err.(*customError.NotFoundError)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.True(t, ok)
		assert.Equal(t, errMessage, err.Error())
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test update success", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		accountRepositoryMock := mck.NewAccountRepository()
		configMock := mckConfig.NewConfigMock()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil).Times(1)

		accountRepositoryMock.Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Account{
				Id:        1,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: helper.StringToDate("2020-10-10 00:00:00"),
				UpdatedAt: helper.StringToDate("2020-10-10 00:00:00"),
			}, nil).Times(1)

		// test
		request := dto.UpdateAccountRequest{
			Id:              1,
			Email:           "reoshby@gmail.com",
			Username:        "rshby",
			Password:        "123456",
			ConfirmPassword: "123456",
		}
		account, err := accountService.Update(context.Background(), &request)

		assert.Nil(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, 1, account.Id)
		assert.Equal(t, "reoshby@gmail.com", account.Email)
		assert.Equal(t, "rshby", account.Username)
		assert.Equal(t, "123456", account.Password)
		assert.Equal(t, "2020-10-10 00:00:00", account.UpdatedAt)
		helperPasswordMock.Mock.AssertExpectations(t)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
}

// unit test method Login
func TestLoginAccountService(t *testing.T) {
	t.Run("test login error validation", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// test
		request := dto.LoginRequest{
			Email:    "reo",
			Password: "123",
		}

		login, err := accountService.Login(context.Background(), &request)
		assert.Nil(t, login)
		assert.NotNil(t, err)
		assert.Error(t, err)

		validationErrors, ok := err.(validator.ValidationErrors)
		assert.True(t, ok)

		for _, errorField := range validationErrors {
			fmt.Println(fmt.Sprintf("%v", errorField.Error()))
		}
	})
	t.Run("test login error account not found", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "record not found"
		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage)).Times(1)

		// test
		request := dto.LoginRequest{
			Email:    "reoshby@gmail.com",
			Password: "123456",
		}

		login, err := accountService.Login(context.Background(), &request)
		assert.Nil(t, login)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())

		_, ok := err.(*customError.NotFoundError)
		assert.True(t, ok)
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test login error password not match", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Account{
				Id:        1,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil).Times(1)

		errMessage := "password not match"
		helperPasswordMock.Mock.On("CheckPasswordHash", mock.Anything, mock.Anything).
			Return(false).Times(1)

		// test
		request := dto.LoginRequest{
			Email:    "reoshby@gmail.com",
			Password: "123456",
		}

		login, err := accountService.Login(context.Background(), &request)
		assert.Nil(t, login)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())

		_, ok := err.(*customError.BadRequestError)
		assert.True(t, ok)

		accountRepositoryMock.Mock.AssertExpectations(t)
		helperPasswordMock.Mock.AssertExpectations(t)
	})
	t.Run("test login success", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		configMock.Mock.On("Config").Return(&config.ConfigApp{
			App:      nil,
			Database: nil,
			Jaeger:   nil,
			Jwt: &config.JWT{
				SecretKey: "sangatrahasia123",
				Issuer:    "coba-metrics-app",
				Subject:   "token",
				Expired:   5,
			},
		}).Times(1)

		accountRepositoryMock.Mock.On("GetByEmail", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Account{
				Id:        0,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil).Times(1)

		helperPasswordMock.Mock.On("CheckPasswordHash", mock.Anything, mock.Anything).
			Return(true).Times(1)

		// test
		request := dto.LoginRequest{
			Email:    "reoshby@gmail.com",
			Password: "123456",
		}

		login, err := accountService.Login(context.Background(), &request)
		assert.Nil(t, err)
		assert.NotNil(t, login)

		loginJson, _ := json.Marshal(&login)
		fmt.Println(string(loginJson))
	})
}

func TestGetAllAccountService(t *testing.T) {
	t.Run("test get all accounts error internal server", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "database refused"
		accountRepositoryMock.Mock.On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage)).Times(1)

		// test
		accounts, err := accountService.GetAll(context.Background(), 10, 2)
		assert.Nil(t, accounts)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test get all accounts error not found", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "record not found"
		accountRepositoryMock.Mock.On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage)).Times(1)

		// test
		all, err := accountService.GetAll(context.Background(), 10, 2)
		assert.Nil(t, all)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
	t.Run("test get all accounts success", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		configMock := mckConfig.NewConfigMock()
		accountRepositoryMock := mck.NewAccountRepository()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountService := service.NewAccountService(db, validate, configMock, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		accountRepositoryMock.Mock.On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Account{
				{
					Id:        1,
					Email:     "reoshby@gmail.com",
					Username:  "rshby",
					Password:  "123456",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					Id:        2,
					Email:     "reoshby2@gmail.com",
					Username:  "rshby2",
					Password:  "123456",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}, nil).Times(1)

		// test
		accounts, err := accountService.GetAll(context.Background(), 2, 1)
		assert.Nil(t, err)
		assert.NotNil(t, accounts)
		assert.Equal(t, 2, len(accounts))
		accountRepositoryMock.Mock.AssertExpectations(t)
	})
}
