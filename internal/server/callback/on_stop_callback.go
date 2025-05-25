package callback

import (
	"context"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/util"
)

func OnServerStoppedCallback(ctx context.Context) error {
	printServerStopped()
	err1 := exportMetricsToPermanentStorage(ctx)
	err2 := closeStorageConnection()
	return util.WrappedErrs(err1, err2)
}

func printServerStopped() {
	fmt.Println("Server shutdown")
}

func exportMetricsToPermanentStorage(ctx context.Context) error {
	if config.Configuration.PermanentStorage == nil {
		fmt.Println("Saving metrics to persistent storage by server shutdowns is disabled")
		return nil
	}

	metricManager := manager.MetricManager{}
	err := metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	(*config.Configuration.PermanentStorage).Close()

	return err
}

func closeStorageConnection() error {
	return config.Configuration.Repository.CloseStorage()
}
