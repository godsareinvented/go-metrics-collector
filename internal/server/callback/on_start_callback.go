package callback

import (
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
)

func OnServerStartedCallback() error {
	printServerStarted()
	return importMetricsFromPermanentStorage()
}

func printServerStarted() {
	fmt.Println("Server is up and running")
}

func importMetricsFromPermanentStorage() error {
	if nil == config.Configuration.PermanentStorage {
		fmt.Println("Saving metrics to and uploading from persistent storage between server outages is disabled")
		return nil
	}

	metricManager := manager.MetricManager{}
	return metricManager.ImportFrom(config.Configuration.PermanentStorage)
}
