package main

import (
	"context"
	"errors"
	clientPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/service/metric"
	"github.com/godsareinvented/go-metrics-collector/internal/general/utils"
	"time"
)

func main() {
	ctx := context.Background()

	configConfigurator := config.ConfigConfigurator{}
	err := configConfigurator.ParseConfig(ctx)
	if nil != err {
		panic(err)
	}

	metricManager, err := metric.NewMetricManager(
		ctx,
		config.Configuration.MetricNameList,
		metric.NewDataCollector(),
		clientPackage.NewClientWithRetry(),
	)
	if nil != err {
		panic(err)
	}

	wg, errCh := metricManager.CollectAndSend()

	for {
		select {
		case errNew, ok := <-errCh:
			err = utils.WrapErrs(err, errNew)
			if !ok {
				err = utils.WrapErrs(errors.New("error channel has been closed"), err)
				panic(err)
			}
		case <-ctx.Done():
			wg.Wait()
			err = utils.WrapErrs(errors.New("context done"), err)
			panic(err)
		}
	}
}
