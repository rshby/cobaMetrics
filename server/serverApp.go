package server

import (
	"cobaMetrics/app/config"
	"cobaMetrics/app/middleware"
	"cobaMetrics/metrics"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type ServerApp struct {
	App  *fiber.App
	Port int
}

func NewServerApp(config config.IConfig) IServer {
	// add metrics
	metrics := metrics.AddMetrics()
	prometheus.MustRegister(metrics.CounterReq, metrics.DurationReq)

	app := fiber.New()

	v1 := app.Group("/api/v1")
	v1.Use(middleware.MetricsMiddleware(config, metrics))

	// add endpoint
	v1.Get("/test", func(ctx *fiber.Ctx) error {
		statusCode := http.StatusOK
		ctx.Status(statusCode)
		return ctx.JSON(&map[string]any{
			"status_code": statusCode,
			"status":      "ok",
			"message":     "success test",
		})
	})

	v1.Get("/users", func(ctx *fiber.Ctx) error {
		statusCode := http.StatusOK
		ctx.Status(statusCode)
		return ctx.JSON(&map[string]any{
			"status_code": statusCode,
			"status":      "ok",
			"message":     "success get users",
		})
	})

	v1.Get("/products", func(ctx *fiber.Ctx) error {
		statusCode := http.StatusOK
		ctx.Status(statusCode)
		return ctx.JSON(&map[string]any{
			"status_code": statusCode,
			"status":      "ok",
			"message":     "success get products",
		})
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return &ServerApp{
		App:  app,
		Port: config.Config().App.Port,
	}
}

func (s *ServerApp) RunServer() error {
	return s.App.Listen(fmt.Sprintf(":%v", s.Port))
}
