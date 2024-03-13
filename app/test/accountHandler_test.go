package test

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/handler"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	mockService "cobaMetrics/app/test/mock/service"
	mockError "cobaMetrics/app/test/mock/validationErrors"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// unit test add account handler
func TestAddAccountHandler(t *testing.T) {
	t.Run("test add account error nil request body", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Post("/", accountHandler.Add)

		// test
		// create request
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		// receive response body
		responseBody := map[string]any{}
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "bad request", responseBody["status"].(string))

		fmt.Println(responseBody["message"].(string))
	})
	t.Run("test add account error error validation", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Post("/", accountHandler.Add)

		// create request body
		req := dto.AddUserRequest{
			Email:    "reo",
			Username: "1",
			Password: "1",
		}
		reqJson, _ := json.Marshal(&req)

		// mock
		accountService.Mock.On("Add", mock.Anything, mock.Anything).
			Return(nil, validator.ValidationErrors{
				&mockError.FieldErrorMock{TagError: "email", FieldErr: "email"},
				&mockError.FieldErrorMock{TagError: "min", FieldErr: "password"},
			})

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		responseBody := map[string]any{}
		body, err := io.ReadAll(response.Body)
		json.Unmarshal(body, &responseBody)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "bad request", responseBody["status"].(string))

		fmt.Println(responseBody["message"].(string))
	})
	t.Run("test add account error internal server error", func(t *testing.T) {
		accountServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Add)

		// create request body
		req := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		reqJson, _ := json.Marshal(&req)

		// mock
		errMessage := "error connection refused"
		accountServiceMock.Mock.On("Add", mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage)).Times(1)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)

		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusInternalServerError, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusInternalServerError), responseBody["status"])
		assert.Equal(t, errMessage, responseBody["message"].(string))
		accountServiceMock.Mock.AssertExpectations(t)
	})
	t.Run("test add user error bad request", func(t *testing.T) {
		accountServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Add)

		// mock
		errMessage := "error bad request"
		accountServiceMock.Mock.On("Add", mock.Anything, mock.Anything).
			Return(nil, customError.NewBadRequestError(errMessage))

		// create request body
		req := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		reqJson, _ := json.Marshal(&req)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
	})
	t.Run("test add user error not found", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Post("/", accountHandler.Add)

		// mock
		errMessage := "error not found"
		accountService.Mock.On("Add", mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage)).Times(1)

		// create request body
		req := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		reqJson, _ := json.Marshal(&req)

		// create request body
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusNotFound), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		accountService.Mock.AssertExpectations(t)

	})
	t.Run("test add user success", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Post("/", accountHandler.Add)

		// mock
		accountService.Mock.On("Add", mock.Anything, mock.Anything).
			Return(&dto.AddUserResponse{
				Id:        1,
				Email:     "reoshby@gmail.com",
				Username:  "rshby",
				CreatedAt: "2020-10-10 00:00:00",
			}, nil).Times(1)

		// create request body
		req := dto.AddUserRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "123456",
		}
		reqJson, _ := json.Marshal(&req)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusOK), responseBody["status"].(string))
		assert.Equal(t, float64(1), responseBody["data"].(map[string]any)["id"].(float64))
		assert.Equal(t, "reoshby@gmail.com", responseBody["data"].(map[string]any)["email"].(string))
		assert.Equal(t, "rshby", responseBody["data"].(map[string]any)["username"].(string))
		assert.Equal(t, "2020-10-10 00:00:00", responseBody["data"].(map[string]any)["created_at"].(string))
	})
}

// unit test get account by email handler
func TestGetByEmailAccountHandler(t *testing.T) {
	t.Run("test get account by email error validation", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Get("/", accountHandler.GetByEmail)

		// mock
		accountService.Mock.On("GetByEmail", mock.Anything, "reo").
			Return(nil, validator.ValidationErrors{
				&mockError.FieldErrorMock{
					TagError: "email",
					FieldErr: "email",
				},
			}).Times(1)

		// create http request
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		r := request.URL.Query()
		r.Add("email", "reo")
		request.URL.RawQuery = r.Encode()

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		accountService.Mock.AssertExpectations(t)
		fmt.Println(responseBody["message"].(string))
	})
	t.Run("test get account by email error internal server", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Get("/", accountHandler.GetByEmail)

		// mock
		email := "reoshby@gmail.com"
		errMessage := "error internal server"
		accountService.Mock.On("GetByEmail", mock.Anything, email).
			Return(nil, customError.NewInternalServerError(errMessage)).Times(1)

		// create http request
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		r := request.URL.Query()
		r.Add("email", email)
		request.URL.RawQuery = r.Encode()

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusInternalServerError, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusInternalServerError), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		accountService.Mock.AssertExpectations(t)
	})
	t.Run("test get account by email error bad request", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Get("/", accountHandler.GetByEmail)

		// mock
		email := "reoshby@gmail.com"
		errMessage := "error bad request"
		accountService.Mock.On("GetByEmail", mock.Anything, email).
			Return(nil, customError.NewBadRequestError(errMessage)).Times(1)

		// create http request
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.Nil(t, err)
		r := request.URL.Query()
		r.Add("email", email)
		request.URL.RawQuery = r.Encode()

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))

		accountService.Mock.AssertExpectations(t)
	})
	t.Run("test get account by email error not found", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Get("/", accountHandler.GetByEmail)

		// mock
		email := "reoshby@gmail.com"
		errMessage := "record not found"
		accountService.Mock.On("GetByEmail", mock.Anything, email).
			Return(nil, customError.NewNotFoundError(errMessage)).Times(1)

		// create http request
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.Nil(t, err)
		assert.NotNil(t, request)

		r := request.URL.Query()
		r.Add("email", email)
		request.URL.RawQuery = r.Encode()

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		// receive response
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusNotFound), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		accountService.Mock.AssertExpectations(t)
	})
	t.Run("test get account by email success", func(t *testing.T) {
		accountService := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountService)

		app := fiber.New()
		app.Get("/", accountHandler.GetByEmail)

		// mock
		email := "reoshby@gmail.com"
		accountService.Mock.On("GetByEmail", mock.Anything, email).
			Return(&dto.AccountDetailResponse{
				Id:        1,
				Email:     email,
				Username:  "rshby",
				Password:  "123456",
				CreatedAt: "2020-10-10 00:00:00",
				UpdatedAt: "2020-10-10 00:00:00",
			}, nil).Times(1)

		// create http test
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NotNil(t, request)
		assert.Nil(t, err)

		r := request.URL.Query()
		r.Add("email", email)
		request.URL.RawQuery = r.Encode()

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusOK), responseBody["status"].(string))
		assert.Equal(t, email, responseBody["data"].(map[string]any)["email"].(string))
		assert.Equal(t, "rshby", responseBody["data"].(map[string]any)["username"].(string))
		assert.Equal(t, "123456", responseBody["data"].(map[string]any)["password"].(string))
	})
}

// unit test handler login
func TestLoginAccountHandler(t *testing.T) {
	t.Run("test login error nil request body", func(t *testing.T) {
		accountServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Login)

		// test
		// create request body
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.NotNil(t, response)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.NotNil(t, body)
		assert.Nil(t, err)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
	})
	t.Run("test login error validation", func(t *testing.T) {
		accountServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Login)

		// mock
		accountServiceMock.Mock.On("Login", mock.Anything, mock.Anything).
			Return(nil, validator.ValidationErrors{
				&mockError.FieldErrorMock{TagError: "email", FieldErr: "email"},
				&mockError.FieldErrorMock{TagError: "min", FieldErr: "password"},
			}).Times(1)

		// create request body
		req := dto.LoginRequest{
			Email:    "reo",
			Password: "123",
		}

		reqJson, _ := json.Marshal(&req)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
		accountServiceMock.Mock.AssertExpectations(t)
	})
	t.Run("test login error bad request", func(t *testing.T) {
		accountServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Login)

		// mock
		errMessage := "error bad request"
		accountServiceMock.Mock.On("Login", mock.Anything, mock.Anything).
			Return(nil, customError.NewBadRequestError(errMessage)).Times(1)

		// test
		// create request body
		req := dto.LoginRequest{
			Email:    "reoshby@gmail.com",
			Password: "123456",
		}

		reqJson, _ := json.Marshal(&req)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusBadRequest), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
	})
	t.Run("test login error not found", func(t *testing.T) {
		accountServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accountServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Login)

		// mock
		errMessage := "record not found"
		accountServiceMock.Mock.On("Login", mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage)).Times(1)

		// test
		// create request body
		req := dto.LoginRequest{
			Email:    "reoshby@gmail.com",
			Password: "123456",
		}

		reqJson, _ := json.Marshal(&req)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.NotNil(t, response)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusNotFound), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		accountServiceMock.Mock.AssertExpectations(t)
	})
	t.Run("test login error internal server", func(t *testing.T) {
		accontServiceMock := mockService.NewAccountServiceMock()
		accountHandler := handler.NewAccountHandler(accontServiceMock)

		app := fiber.New()
		app.Post("/", accountHandler.Login)

		// mock
		errMessage := "error internal server"
		accontServiceMock.Mock.On("Login", mock.Anything, mock.Anything).
			Return(nil, customError.NewInternalServerError(errMessage)).Times(1)

		// test
		// create request body
		req := dto.LoginRequest{
			Email:    "reoshby@gmail.com",
			Password: "123456",
		}

		reqJson, _ := json.Marshal(&req)

		// create http request
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqJson)))
		request.Header.Add("Content-Type", "application/json")

		// receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

		// receive response body
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.NotNil(t, body)

		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusInternalServerError, int(responseBody["status_code"].(float64)))
		assert.Equal(t, helper.CodeToStatus(http.StatusInternalServerError), responseBody["status"].(string))
		assert.Equal(t, errMessage, responseBody["message"].(string))
		accontServiceMock.Mock.AssertExpectations(t)
	})
}
