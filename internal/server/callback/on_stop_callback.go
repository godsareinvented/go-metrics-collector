package callback

import (
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
)

func OnServerStoppedCallback() error {
	var resultError error

	printServerStopped()

	err := exportMetricsToPermanentStorage()
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

func exportMetricsToPermanentStorage() error {
	if nil == config.Configuration.PermanentStorage {
		fmt.Println("Saving metrics to persistent storage by server shutdowns is disabled")
		return nil
	}

	metricManager := manager.MetricManager{}
	err := metricManager.ExportTo(config.Configuration.PermanentStorage)
	(*config.Configuration.PermanentStorage).Close()

	return err
}

func closeStorage() error {
	return config.Configuration.Repository.CloseStorage()
}
