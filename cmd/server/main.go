package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	webServer := server.Server{}
	webServer.Start()
}
