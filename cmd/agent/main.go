package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/client"
	"github.com/oldhanasong/go-metrics-collector/internal/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	c := client.NewClientWithRetry()
	c.Use(decorator.GzipCompress)

	metricManager := metric.MetricManager{
		MetricList:          dictionary.MetricNameList[:],
		MetricDataCollector: &data_collector.MetricDataCollector{},
		Client:              c,
	}
	metricManager.Init()

	metricManager.CollectAndSend(context.Background())
	select {}
}
