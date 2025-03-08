package main

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/server"
	"github.com/godsareinvented/go-metrics-collector/internal/server/server/callback"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	configConfigurator := config.ConfigConfigurator{}
	err := configConfigurator.GetConfig(ctx)
	if nil != err {
		panic(err)
	}

	webServer := server.Server{
		OnStart: callback.OnServerStartedCallback,
		OnStop:  callback.OnServerStoppedCallback,
	}
	go webServer.Start()

	select {}
}
