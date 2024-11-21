package callback

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
)

func OnServerStoppedCallback() error {
	printServerStopped()
	return exportMetricsToPermanentStorage()
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
