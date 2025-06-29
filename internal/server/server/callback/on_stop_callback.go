package callback

import (
	"context"
	"errors"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
	manager "github.com/oldhanasong/go-metrics-collector/internal/server/service/metric"
	"go.uber.org/multierr"
)

func OnServerStoppedCallback(ctx context.Context) error {
	printServerStopped()
	err1 := exportMetricsToPermanentStorage(ctx)
	err2 := closeStorageConnection()
	return multierr.Combine(err1, err2)
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
	err := config.Configuration.Repository.CloseStorage()
	if !errors.Is(err, repository.ErrDontImplementConnectorInterface) {
		return err
	}
	return nil
}
