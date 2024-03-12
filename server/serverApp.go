package server

import (
	"cobaMetrics/app/config"
	"cobaMetrics/app/handler"
	"cobaMetrics/app/helper"
	"cobaMetrics/app/middleware"
	"cobaMetrics/app/repository"
	"cobaMetrics/app/service"
	"cobaMetrics/metrics"
	"cobaMetrics/router"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServerApp struct {
	App  *fiber.App
	Port int
}

func NewServerApp(config config.IConfig, db *sql.DB, validate *validator.Validate) IServer {
	// add metrics
	metrics := metrics.AddMetrics()
	prometheus.MustRegister(metrics.CounterReq, metrics.DurationReq)

	// register repository
	accountRepository := repository.NewAccountRepository()

	// register service
	accountService := service.NewAccountService(db, validate, accountRepository, helper.NewHelperPassword())

	// register handler
	accountHandler := handler.NewAccountHandler(accountService)

	// create instance fiber
	app := fiber.New()

	v1 := app.Group("/api/v1")
	v1.Use(middleware.MetricsMiddleware(config, metrics))

	// router
	router.GenerateAccountRouter(v1, accountHandler)

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return &ServerApp{
		App:  app,
		Port: config.Config().App.Port,
	}
}

func (s *ServerApp) RunServer() error {
	return s.App.Listen(fmt.Sprintf(":%v", s.Port))
}
