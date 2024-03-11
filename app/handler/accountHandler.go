package handler

import (
	"cobaMetrics/app/customError"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	IService "cobaMetrics/app/service/interface"
	"encoding/json"
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
