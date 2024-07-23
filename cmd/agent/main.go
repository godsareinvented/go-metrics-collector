package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
)

func main() {
	counterMetricManager := manager.MetricManager[int64]{
		MetricList:          dictionary.CounterMetricNameList[:],
		MetricDataCollector: &data_collector.CounterMetricDataCollector{},
	}

	gaugeMetricManager := manager.MetricManager[float64]{
		MetricList:          dictionary.GaugeMetricNameList[:],
		MetricDataCollector: &data_collector.GaugeMetricDataCollector{},
	}

	go counterMetricManager.CollectAndSend(context.Background())
	go gaugeMetricManager.CollectAndSend(context.Background())
	select {}
}
