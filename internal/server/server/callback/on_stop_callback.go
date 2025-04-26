package callback

import (
	"context"
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric"
)

func OnServerStoppedCallback(ctx context.Context, errCh chan<- error) {
	printServerStopped()

	err := exportMetricsToPermanentStorage(ctx)
	if nil != err {
		errCh <- err
		return
	}

	err = closeResources()
	if nil != err {
		errCh <- err
		return
	}

	return
}

func printServerStopped() {
	fmt.Println("Server shutdown")
}

func exportMetricsToPermanentStorage(ctx context.Context) error {
	if nil == config.Configuration.PermanentStorage {
		fmt.Println("Saving metrics to persistent storage by server shutdowns is disabled")
		return nil
	}

	metricManager := metric.MetricManager{}
	return metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
}

func closeResources() error {
	(*config.Configuration.PermanentStorage).Close()
	return config.Configuration.Repository.CloseStorage()
}
