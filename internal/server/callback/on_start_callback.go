package callback

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"time"
)

func OnServerStartedCallback() error {
	printServerStarted()
	initExportTask()
	return importMetricsFromPermanentStorage()
}

func printServerStarted() {
	fmt.Println("Server is up and running")
}

func initExportTask() {
	if config.Configuration.StoreInterval > 0 {
		go exportTask()
	}
}

func exportTask() {
	metricManager := manager.MetricManager{}
	ticker := time.NewTicker(time.Duration(config.Configuration.StoreInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = metricManager.ExportTo(config.Configuration.PermanentStorage)
		}
	}
}

func importMetricsFromPermanentStorage() error {
	if config.Configuration.PermanentStorage == nil {
		fmt.Println("Saving metrics to and uploading from persistent storage between server outages is disabled")
		return nil
	}

	if !config.Configuration.Restore {
		return nil
	}

	metricManager := manager.MetricManager{}
	return metricManager.ImportFrom(config.Configuration.PermanentStorage)
}
