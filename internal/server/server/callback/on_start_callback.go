package callback

import (
	"context"
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric"
	"time"
)

func OnServerStartedCallback(ctx context.Context, errCh chan<- error) error {
	printServerStarted()

	err := importMetricsFromPermanentStorage(ctx)
	if nil != err {
		return err
	}

	err = initTasks(ctx, errCh)
	if nil != err {
		return err
	}

	return nil
}

func printServerStarted() {
	fmt.Println("Server is up and running")
}

// todo: В будущем обязательно переписать на более надёжную схему.
func initTasks(ctx context.Context, errCh chan<- error) error {
	if config.Configuration.StoreInterval > 0 {
		err := exportTask(ctx, errCh)
		if nil != err {
			return err
		}
	}
	return nil
}

func exportTask(ctx context.Context, errCh chan<- error) error {
	metricManager := metric.MetricManager{}
	err := metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	if nil != err {
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Duration(config.Configuration.StoreInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
				if nil != err {
					errCh <- err
				}
			}
		}
	}()

	return nil
}

func importMetricsFromPermanentStorage(ctx context.Context) error {
	if nil == config.Configuration.PermanentStorage {
		fmt.Println("Saving metrics to and uploading from persistent storage between server outages is disabled")
		return nil
	}

	metricManager := metric.MetricManager{}
	return metricManager.ImportFrom(ctx, config.Configuration.PermanentStorage)
}
