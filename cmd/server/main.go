package main

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/server"
	"github.com/oldhanasong/go-metrics-collector/internal/server/server/callback"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
		<-exitCh
		cancel()
	}()

	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	webServer := server.Server{
		OnStart: callback.OnServerStartedCallback,
		OnStop:  callback.OnServerStoppedCallback,
	}
	webServer.Start(ctx)

	select {
	case <-ctx.Done():
		if err := webServer.Stop(); err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}
}
