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

	int64MetricManager := manager.MetricManager[int64]{
		MetricList:          dictionary.Int64MetricNameList[:],
		MetricDataCollector: &metric_data_collector.Int64MetricDataCollector{},
	}

	float64MetricManager := manager.MetricManager[float64]{
		MetricList:          dictionary.Float64MetricNameList[:],
		MetricDataCollector: &metric_data_collector.Float64MetricDataCollector{},
	}

	var n time.Duration = 2
	for {
		int64MetricManager.CollectAndSend()
		float64MetricManager.CollectAndSend()

		time.Sleep(n * time.Second)
	}
}
