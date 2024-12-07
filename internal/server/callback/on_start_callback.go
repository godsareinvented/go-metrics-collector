package callback

import (
	"context"
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"time"
)

func OnServerStartedCallback(ctx context.Context) error {
	printServerStarted()
	initExportTask(ctx)
	return importMetricsFromPermanentStorage(ctx)
}

func printServerStarted() {
	fmt.Println("Server is up and running")
}

// todo: В будущем обязательно переписать на более надёжную схему.
func initExportTask(ctx context.Context) {
	if config.Configuration.StoreInterval > 0 {
		go exportTask(ctx)
	}
}

func exportTask(ctx context.Context) {
	metricManager := manager.MetricManager{}
	ticker := time.NewTicker(time.Duration(config.Configuration.StoreInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
		}
	}
}

func importMetricsFromPermanentStorage(ctx context.Context) error {
	if nil == config.Configuration.PermanentStorage {
		fmt.Println("Saving metrics to and uploading from persistent storage between server outages is disabled")
		return nil
	}

	metricManager := manager.MetricManager{}
	return metricManager.ImportFrom(ctx, config.Configuration.PermanentStorage)
}
