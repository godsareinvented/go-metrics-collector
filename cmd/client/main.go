package main

import (
	"context"
	clientPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/service/metric"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/threading_pattern"
	"time"
)

// todo: Добавить в будущем аналогично серверу контекст.
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	configConfigurator := config.ConfigConfigurator{}
	err := configConfigurator.ParseConfig()
	if nil != err {
		panic(err)
	}

	client := clientPackage.NewClientWithRetry()
	metricManager, err := metric.NewMetricManager(
		metric.NewDataCollector(),
		config.Configuration.MetricNameList,
	)
	if nil != err {
		panic(err)
	}

	errCh := make(chan error)
	defer close(errCh)

	ch := CollectMetrics(ctx, errCh, &metricManager)
	SendMetrics(ch, client)

	select {
	case err := <-errCh:
		if nil != err {
			cancel()
			close(errCh)
			panic(err)
		}
	}
}

func CollectMetrics(ctx context.Context, errCh chan<- error, metricManager *metric.MetricManager) chan *[]dto.Metrics {
	ch := make(chan *[]dto.Metrics)

	go func() {
		for {
			metricList, err := metricManager.Collect(ctx)
			if nil != err {
				close(ch)
				errCh <- err
				return
			}
			ch <- metricList

			if config.Configuration.PollInterval > 0 {
				time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
			}
		}
	}()

	return ch
}

func SendMetrics(inputCh <-chan *[]dto.Metrics, client interfaces.Client) {
	ch := make(chan *[]dto.Metrics)
	_ = threading_pattern.InitWorkerPool(config.Configuration.RateLimit, ch, func(_ int, metricList *[]dto.Metrics) {
		_ = client.SendBatch(metricList)
	})

	go func() {
		defer close(ch)
		for metricList := range inputCh {
			ch <- metricList

			if config.Configuration.PollInterval > 0 {
				time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
			}
		}
	}()
}
