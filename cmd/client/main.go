package main

import (
	"container/list"
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

	var metricQueue = list.New()
	client := clientPackage.NewClientWithRetry()
	metricManager := manager.MetricManager{
		MetricList:    dictionary.MetricNameList[:],
		DataCollector: &metricdatacollector.MetricDataCollector{},
	}
	metricManager.Init()

	go CollectMetrics(metricQueue, &metricManager)
	go SendMetrics(metricQueue, client)

	select {}
}

func CollectMetrics(metricQueue *list.List, metricManager *manager.MetricManager) {
	for {
		// todo: Обернуть элемент очереди в какую-то iterable структуру? Типа, пакет метрик..
		metricQueue.PushBack(metricManager.Collect())

		time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
	}
}

func SendMetrics(metricQueue *list.List, client interfaces.Client) {
	for {
		metricListRaw := metricQueue.Front()
		if nil == metricListRaw {
			time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
			continue
		}

		metricList, ok := metricListRaw.Value.([]dto.Metrics)
		if !ok {
			panic("SendMetrics: metricListRaw is not []dto.Metrics")
		}
		metricQueue.Remove(metricListRaw)

		if nil != metricList {
			_ = client.SendBatch(metricList)
		}

		time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
	}
}
