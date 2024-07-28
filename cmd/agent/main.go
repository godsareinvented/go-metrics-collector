package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
)

func main() {
	metricManager := manager.MetricManager{
		MetricList:          dictionary.MetricNameList[:],
		MetricDataCollector: &data_collector.MetricDataCollector{},
	}

	go metricManager.CollectAndSend(context.Background())
	select {}
}
