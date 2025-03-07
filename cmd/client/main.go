package main

import (
	"context"
	"errors"
	clientPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/service/metric"
)

// todo: Добавить в будущем аналогично серверу контекст.
func main() {
	ctx := context.Background()

	configConfigurator := config.ConfigConfigurator{}
	err := configConfigurator.ParseConfig(ctx)
	if nil != err {
		panic(err)
	}

	metricManager, err := metric.NewMetricManager(
		metric.NewDataCollector(),
		clientPackage.NewClientWithRetry(),
		config.Configuration.MetricNameList,
	)
	if nil != err {
		panic(err)
	}

	errCh := metricManager.CollectAndSend(ctx)

	select {
	case err = <-errCh:
		if nil != err {
			panic(err)
		}
		panic(errors.New("unexpected closure of the error channel"))
	}
}
