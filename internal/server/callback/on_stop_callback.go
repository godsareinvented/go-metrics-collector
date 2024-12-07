package callback

import (
	"context"
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
)

func OnServerStoppedCallback(ctx context.Context) error {
	var resultError error

	printServerStopped()

	err := exportMetricsToPermanentStorage(ctx)
	if nil != err {
		resultError = err
	}

	err = closeStorage()
	if nil != err {
		resultError = err
	}

	return resultError
}

func printServerStopped() {
	fmt.Println("Server shutdown")
}

func exportMetricsToPermanentStorage(ctx context.Context) error {
	if nil == config.Configuration.PermanentStorage {
		fmt.Println("Saving metrics to persistent storage by server shutdowns is disabled")
		return nil
	}

	metricManager := manager.MetricManager{}
	err := metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	(*config.Configuration.PermanentStorage).Close()

	return err
}

func closeStorage() error {
	return config.Configuration.Repository.CloseStorage()
}
