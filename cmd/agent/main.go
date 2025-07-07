package main

import (
	"context"
	"errors"
	blconfig "github.com/oldhanasong/go-metrics-collector/internal/agent/business_logic/config"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/business_logic/metric/data_collector"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric"
	"os"
	"os/signal"
	"sync"
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
	if err := configConfigurator.ParseConfig(); err != nil {
		panic(err)
	}

	c := client.NewClientWithRetry()
	c.Use(decorator.GzipCompress)
	c.Use(decorator.HashCalculation)

	metricManager, err := manager.New(blconfig.Configuration.MetricsToCollect, data_collector.New(), c)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	metricManager.CollectAndSend(ctx, func(err error) {
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
		wg.Done()
	})

	select {
	case <-ctx.Done():
		wg.Wait()
	}
}
