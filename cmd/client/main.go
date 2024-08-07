package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/client"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/data_collector"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"time"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	var metricDTOList []dto.Metric
	metricSender := client.NewInstance()
	metricManager := manager.MetricManager{
		MetricList:    dictionary.MetricNameList[:],
		DataCollector: &metric_data_collector.MetricDataCollector{},
	}
	metricManager.Init()

	go CollectMetrics(&metricDTOList, &metricManager)
	go SendMetrics(&metricDTOList, metricSender)

	select {}
}

func CollectMetrics(metricDTOList *[]dto.Metric, metricManager *manager.MetricManager) {
	for {
		*metricDTOList = metricManager.Collect()

		time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
	}
}

func SendMetrics(metricDTOList *[]dto.Metric, client *client.MetricSender) {
	for {
		for _, metricDTO := range *metricDTOList {
			_ = client.Send(metricDTO)
		}

		time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
	}
}
