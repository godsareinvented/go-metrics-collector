package main

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
)

func main() {
	counterMetricRepository := repository.NewInstance[int64](mem_storage.NewInstance())
	counterMetricManager := manager.MetricManager[int64]{
		MetricList:          dictionary.CounterMetricNameList[:],
		MetricDataCollector: &data_collector.CounterMetricDataCollector{},
		Repository:          counterMetricRepository,
	}

	gaugeMetricRepository := repository.NewInstance[float64](mem_storage.NewInstance())
	gaugeMetricManager := manager.MetricManager[float64]{
		MetricList:          dictionary.GaugeMetricNameList[:],
		MetricDataCollector: &data_collector.GaugeMetricDataCollector{},
		Repository:          gaugeMetricRepository,
	}

	go counterMetricManager.CollectAndSend(context.Background())
	go gaugeMetricManager.CollectAndSend(context.Background())
	select {}
}
