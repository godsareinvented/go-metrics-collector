package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/data_collector"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
	"time"
)

func main() {
	memStorage := mem_storage.NewInstance()
	repository.NewInstance(memStorage)

	metricManager := manager.MetricManager{
		MetricList:          dictionary.MetricNameList[:],
		MetricDataCollector: &metric_data_collector.MetricDataCollector{},
	}

	var n time.Duration = 2
	for {
		metricManager.CollectAndSend()

		time.Sleep(n * time.Second)
	}
}
