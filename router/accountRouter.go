package router

import (
	"cobaMetrics/app/handler"
	"github.com/gofiber/fiber/v2"
)

func GenerateAccountRouter(app fiber.Router, handler *handler.AccountHandler) {
	app.Post("/account", handler.Add)
	app.Get("/account", handler.GetByEmail)
}
