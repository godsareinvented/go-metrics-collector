package main

import (
	clientPackage "github.com/godsareinvented/go-metrics-collector/internal/client"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	metricdatacollector "github.com/godsareinvented/go-metrics-collector/internal/service/metric/data_collector"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"time"
)

// todo: Добавить в будущем аналогично серверу контекст.
func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	var metricDTOList []dto.Metrics
	client := clientPackage.NewClientWithRetry()
	metricManager := manager.MetricManager{
		MetricList:    dictionary.MetricNameList[:],
		DataCollector: &metricdatacollector.MetricDataCollector{},
	}
	metricManager.Init()

	go CollectMetrics(&metricDTOList, &metricManager)
	go SendMetrics(&metricDTOList, client)

	select {}
}

func CollectMetrics(metricDTOList *[]dto.Metrics, metricManager *manager.MetricManager) {
	for {
		*metricDTOList = metricManager.Collect()

		time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
	}
}

func SendMetrics(metricDTOList *[]dto.Metrics, client interfaces.Client) {
	for {
		if nil != *metricDTOList {
			_ = client.SendBatch(*metricDTOList)
		}

		time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
	}
}
