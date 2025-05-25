package callback

import (
	"context"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
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
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
		}
	}
}

func importMetricsFromPermanentStorage(ctx context.Context) error {
	if config.Configuration.PermanentStorage == nil {
		fmt.Println("Saving metrics to and uploading from persistent storage between server outages is disabled")
		return nil
	}

	if !config.Configuration.Restore {
		return nil
	}

	metricManager := manager.MetricManager{}
	return metricManager.ImportFrom(ctx, config.Configuration.PermanentStorage)
}
