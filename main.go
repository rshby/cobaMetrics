package main

import (
	config "cobaMetrics/app/config"
	"cobaMetrics/app/tracing"
	"cobaMetrics/database"
	"cobaMetrics/server"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// load config
	config := config.NewConfigApp()

	// start tracing
	tracer, closer := tracing.NewJaegerTracing(config)
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	// connect db
	db := database.ConnectDB(config)

	validate := validator.New()

	// run server
	server := server.NewServerApp(config, db, validate)

	server.RunServer()
}
