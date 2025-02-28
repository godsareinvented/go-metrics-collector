package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric/data_collector"
)

func main() {
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

	metricManager.CollectAndSend(context.Background(), func(err error) {
		panic(err)
	})

	select {}
}
