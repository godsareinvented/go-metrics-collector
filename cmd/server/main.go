package main

import (
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server"
	"github.com/oldhanasong/go-metrics-collector/internal/server/callback"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	webServer := server.Server{
		OnStart: callback.OnServerStartedCallback,
		OnStop:  callback.OnServerStoppedCallback,
	}
	webServer.Start()

	select {}
}
