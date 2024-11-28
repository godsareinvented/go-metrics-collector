package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server"
	"github.com/oldhanasong/go-metrics-collector/internal/server/callback"
)

func main() {
	ctx := context.Background()

	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	webServer := server.Server{
		OnStart: callback.OnServerStartedCallback,
		OnStop:  callback.OnServerStoppedCallback,
	}
	webServer.Start(ctx)

	select {}
}
