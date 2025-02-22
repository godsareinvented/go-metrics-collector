package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric/data_collector"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	c := client.NewClientWithRetry()
	c.Use(decorator.GzipCompress)
	c.Use(decorator.HashCalculation)

	metricManager, err := manager.New(dictionary.MetricNameList[:], data_collector.New(), c)
	if err != nil {
		panic(err)
	}

	err = metricManager.CollectAndSend(context.Background())
	if err != nil {
		panic(err)
	}

	select {}
}
