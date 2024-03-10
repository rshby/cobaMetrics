package main

import (
	config "cobaMetrics/app/config"
	"cobaMetrics/server"
)

func main() {

	// load config
	config := config.NewConfigApp()

	// run server
	server := server.NewServerApp(config)

	server.RunServer()
}
