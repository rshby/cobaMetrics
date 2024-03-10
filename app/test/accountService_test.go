package test

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/model/dto"
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
}
