package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
)

func main() {
	metricManager := manager.MetricManager{MetricDataCollector: data_collector.MetricDataCollector{}}
	metricManager.CollectAndSend(context.Background())
}
