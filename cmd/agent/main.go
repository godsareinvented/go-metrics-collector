package main

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric/data_collector"
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

	metricManager, err := manager.New(config.Configuration.MetricsToCollect, data_collector.New(), c)
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
