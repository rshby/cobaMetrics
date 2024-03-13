package router

import (
	"cobaMetrics/app/handler"
	"github.com/gofiber/fiber/v2"
)

func GenerateAccountRouter(app fiber.Router, authMiddleware fiber.Handler, handler *handler.AccountHandler) {
	app.Post("/account", handler.Add)
	app.Get("/account", authMiddleware, handler.GetByEmail)
	app.Post("/login", handler.Login)
}
