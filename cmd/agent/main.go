package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/metric/data_collector"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	c := client.NewClientWithRetry()
	c.Use(decorator.GzipCompress)
	c.Use(decorator.HashCalculation)

	metricManager := metric.MetricManager{
		MetricList:          dictionary.MetricNameList[:],
		MetricDataCollector: &data_collector.MetricDataCollector{},
		Client:              c,
	}
	metricManager.Init()

	metricManager.CollectAndSend(context.Background())
	select {}
}
