package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/data_collector"
	"time"
)

func main() {
	var n time.Duration = 2

	for {
		metricManager := manager.MetricManager{MetricDataCollector: metric_data_collector.MetricDataCollector{}}
		metricManager.CollectAndSend()

		time.Sleep(n * time.Second)
	}
}
