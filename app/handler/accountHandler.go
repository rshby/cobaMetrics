package handler

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	IService "cobaMetrics/app/service/interface"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"strings"
)

type AccountHandler struct {
	AccountService IService.IAccountService
}

func NewAccountHandler(accService IService.IAccountService) *AccountHandler {
	return &AccountHandler{accService}
}

// handler insert
func (a *AccountHandler) Add(ctx *fiber.Ctx) error {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx.Context(), "AccountHandler Add")
	defer span.Finish()

	// decode request body
	var request dto.AddUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		ext.Error.Set(span, true)

		statusCode := http.StatusBadRequest
		ctx.Status(statusCode)
		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}

		// log response with tracing
		responseJson, _ := json.Marshal(&response)
		span.LogFields(
			log.String("response", string(responseJson)))

		return ctx.JSON(&response)
	}

	// log request with tracing
	requestJson, _ := json.Marshal(&request)
	span.LogFields(
		log.String("request", string(requestJson)),
	)

	// call procedure in service
	account, err := a.AccountService.Add(ctxTracing, &request)
	if err != nil {
		ext.Error.Set(span, true)

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errMessage := []string{}
			for _, fieldError := range validationErrors {
				errMessage = append(errMessage, fieldError.Error())
			}

			statusCode := http.StatusBadRequest
			ctx.Status(statusCode)
			response := dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.CodeToStatus(statusCode),
				Message:    strings.Join(errMessage, ". "),
			}

			// log with tracing
			responseJson, _ := json.Marshal(&response)
			span.LogFields(log.String("response", string(responseJson)))
			return ctx.JSON(&response)
		}

		var statusCode int
		switch err.(type) {
		case *customError.InternalServerError:
			statusCode = http.StatusInternalServerError
		case *customError.BadRequestError:
			statusCode = http.StatusBadRequest
		case *customError.NotFoundError:
			statusCode = http.StatusNotFound
		}

		ctx.Status(statusCode)

		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}

		// log response with tracing
		responseJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(responseJson)))

		return ctx.JSON(&response)
	}

	// sucess
	statusCode := http.StatusOK
	response := dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.CodeToStatus(statusCode),
		Message:    "success add new account",
		Data:       account,
	}

	// log with tracing
	responseJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(responseJson)))

	return ctx.JSON(&response)
}

// handler get data accounts by email
func (a *AccountHandler) GetByEmail(ctx *fiber.Ctx) error {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx.Context(), "AccountHandler GetByEmail")
	defer span.Finish()

	// get from params
	email := ctx.Query("email")

	// call procedure in service
	account, err := a.AccountService.GetByEmail(ctxTracing, email)
	if err != nil {
		ext.Error.Set(span, true)

		// error validation
		if _, ok := err.(validator.ValidationErrors); ok {
			statusCode := http.StatusBadRequest
			response := dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.CodeToStatus(statusCode),
				Message:    err.Error(),
			}

			// log with tracing
			resJson, _ := json.Marshal(&response)
			span.LogFields(log.String("response", string(resJson)))

			ctx.Status(statusCode)
			return ctx.JSON(&response)
		}

		var statusCode int
		switch err.(type) {
		case *customError.BadRequestError:
			statusCode = http.StatusBadRequest
		case *customError.NotFoundError:
			statusCode = http.StatusNotFound
		case *customError.InternalServerError:
			statusCode = http.StatusInternalServerError
		}

		// create response object
		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}

		// log with tracing
		responseJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(responseJson)))

		ctx.Status(statusCode)
		return ctx.JSON(&response)
	}

	// success
	statusCode := http.StatusOK

	// create response object
	response := dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.CodeToStatus(statusCode),
		Message:    "success get data account",
		Data:       account,
	}

	// log with tracing
	responseJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(responseJson)))

	ctx.Status(statusCode)
	return ctx.JSON(&response)
}

// handler login
func (a *AccountHandler) Login(ctx *fiber.Ctx) error {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx.Context(), "AccountHandler Login")
	defer span.Finish()

	// decode request_body
	var request dto.LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		statusCode := http.StatusBadRequest
		ctx.Status(statusCode)
		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}

		// log with tracing
		ext.Error.Set(span, true)
		responseJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(responseJson)))

		return ctx.JSON(&response)
	}

	// log request with tracing
	requestJson, _ := json.Marshal(&request)
	span.LogFields(log.String("request", string(requestJson)))

	// call procedure in service
	login, err := a.AccountService.Login(ctxTracing, &request)
	if err != nil {
		ext.Error.Set(span, true)

		// error validation
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMessage []string
			for _, errorField := range validationErrors {
				errMessage = append(errMessage, fmt.Sprintf("error on field [%v] with tag [%v]", errorField.Field(), errorField.Tag()))
			}

			statusCode := http.StatusBadRequest
			ctx.Status(statusCode)
			response := dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.CodeToStatus(statusCode),
				Message:    strings.Join(errMessage, ". "),
			}

			// log with tracing
			responseJson, _ := json.Marshal(&response)
			span.LogFields(log.String("response", string(responseJson)))

			return ctx.JSON(&response)
		}

		var statusCode int
		switch err.(type) {
		case *customError.BadRequestError:
			statusCode = http.StatusBadRequest
		case *customError.NotFoundError:
			statusCode = http.StatusNotFound
		case *customError.InternalServerError:
			statusCode = http.StatusInternalServerError
		}

		ctx.Status(statusCode)
		response := dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.CodeToStatus(statusCode),
			Message:    err.Error(),
		}

		// log with tracing
		responseJson, _ := json.Marshal(&response)
		span.LogFields(log.String("response", string(responseJson)))

		return ctx.JSON(&response)
	}

	// success login
	statusCode := http.StatusOK
	ctx.Status(statusCode)
	response := dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.CodeToStatus(statusCode),
		Message:    "success login",
		Data:       login,
	}

	// log response with tracing
	responseJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(responseJson)))

	return ctx.JSON(&response)
}
