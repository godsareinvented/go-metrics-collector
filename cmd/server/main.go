package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server"
	"github.com/godsareinvented/go-metrics-collector/internal/server/callback"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	webServer := server.Server{
		OnStart: callback.OnServerStartedCallback,
		OnStop:  callback.OnServerStoppedCallback,
	}
	// todo: Нужен контекст для завершения этой горутины?..
	go webServer.Start()

	select {}
}
