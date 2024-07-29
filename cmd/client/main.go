package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/data_collector"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/sender"
	"time"
)

func main() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	var metricDTOList []dto.Metric
	metricSender := sender.NewSender()
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

func SendMetrics(metricDTOList *[]dto.Metric, metricSender *sender.MetricSender) {
	for {
		for _, metricDTO := range *metricDTOList {
			metricSender.Send(metricDTO)
		}

		time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
	}
}
