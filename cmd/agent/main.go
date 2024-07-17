package main

import (
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
	"time"
)

func main() {
	var n time.Duration = 2

	for {
		metricManager := manager.MetricManager{MetricDataCollector: data_collector.MetricDataCollector{}}
		metricManager.CollectAndSend()

		time.Sleep(n * time.Second)
	}
}
