package test

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	"cobaMetrics/app/model/entity"
	"cobaMetrics/app/service"
	mckHelper "cobaMetrics/app/test/mock/helper"
	mck "cobaMetrics/app/test/mock/repository"
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddUserService(t *testing.T) {
	t.Run("add account error validate", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountRepository := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, accountRepository, helperPasswordMock)

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
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, accountRepositoryMock, helperPasswordMock)

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
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, accountRepositoryMock, helperPassword)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPassword.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

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
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, accountRepositoryMock, helperPassword)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPassword.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

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
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

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
	t.Run("add account success", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.Nil(t, err)

		validate := validator.New()
		helperPasswordMock := mckHelper.NewHelperPasswordMock()
		accountRepositoryMock := mck.NewAccountRepository()
		accountService := service.NewAccountService(db, validate, accountRepositoryMock, helperPasswordMock)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		helperPasswordMock.Mock.On("HashPassword", mock.Anything).
			Return("123456", nil)

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