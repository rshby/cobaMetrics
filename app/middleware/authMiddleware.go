package middleware

import (
	"cobaMetrics/app/config"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/model/dto"
	jwtModel "cobaMetrics/app/model/jwt"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	logFiber "github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"strings"
)

func AuthMiddleware(config config.IConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// create span tracing
		span, _ := opentracing.StartSpanFromContext(ctx.Context(), "Middleware Auth")
		defer span.Finish()

		// log
		request := ctx.Request().Body()
		span.LogFields(log.String("request", string(request)))

		jwtConfig := config.Config().Jwt
		logFiber.Info("masuk middleware auth")

		// get token
		tokenHeader := ctx.Get("authorization")

		if tokenHeader == "" {
			statusCode := http.StatusUnauthorized
			response := dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.CodeToStatus(statusCode),
				Message:    "token required",
			}

			responseJson, _ := json.Marshal(&response)
			span.LogFields(log.String("response", string(responseJson)))
			ext.Error.Set(span, true)

			ctx.Status(statusCode)
			return ctx.JSON(&response)
		}

		tokenString := strings.Split(tokenHeader, " ")
		var token string = tokenString[1]

		// decode claims
		var claims jwtModel.Claims
		tokenWithClaims, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.SecretKey), nil
		})

		claimsJson, _ := json.Marshal(&claims)
		span.LogFields(log.String("claims", string(claimsJson)))

		// if token not valid
		if err != nil {
			ext.Error.Set(span, true)

			statuscCode := http.StatusUnauthorized

			response := dto.ApiResponse{
				StatusCode: statuscCode,
				Status:     helper.CodeToStatus(statuscCode),
				Message:    err.Error(),
			}

			responseJson, _ := json.Marshal(&response)
			span.LogFields(log.String("response", string(responseJson)))

			ctx.Status(statuscCode)
			return ctx.JSON(&response)
		}

		// if not valid at all
		if !tokenWithClaims.Valid {
			statusCode := http.StatusUnauthorized
			response := dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.CodeToStatus(statusCode),
				Message:    "token not valid",
			}

			responseJson, _ := json.Marshal(&response)
			span.LogFields(log.String("response", string(responseJson)))
			ext.Error.Set(span, true)

			ctx.Status(statusCode)
			return ctx.JSON(&response)
		}

		// lolos semua validasi auth
		ctx.Next()
		return nil
	}
}
