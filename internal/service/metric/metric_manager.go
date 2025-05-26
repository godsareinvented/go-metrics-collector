package metric

import (
	"context"
	parserAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/parser"
	valueHandlerAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/value_handler"
	"github.com/oldhanasong/go-metrics-collector/internal/client"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"go.uber.org/multierr"
	"time"
)

type MetricManager struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
	strategies          map[string]interfaces.ParsingStrategy
	client              *client.MetricSender
}

var (
	metricList []dto.Metrics
)

func (metricManager *MetricManager) CollectAndSend(ctx context.Context) {
	if metricManager.MetricList == nil {
		panic("metric list is empty")
	}

	go metricManager.collect(ctx)
	go metricManager.send(ctx)

	select {
	case <-ctx.Done():
		return
	}
}

func (metricManager *MetricManager) collect(ctx context.Context) {
	var metricCollectedData dto.CollectedMetricData
	var metrics []dto.Metrics
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricManager.MetricDataCollector.CollectMetricData(&metricCollectedData)

			metrics = []dto.Metrics{}
			for _, metricName := range metricManager.MetricList {
				strategy := metricManager.strategies[metricName]
				m := strategy.GetMetric(metricName, metricCollectedData)
				metrics = append(metrics, m)
			}
			metricList = metrics

			time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
		}
	}
}

func (metricManager *MetricManager) send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, metrics := range metricList {
				_ = metricManager.client.Send(metrics)
			}

			time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
		}
	}
}

func (metricManager *MetricManager) UpdateMetrics(ctx context.Context, metric dto.Metrics) error {
	repos := config.Configuration.Repository

	metricFromStorage, isSet, err := repos.GetMetric(ctx, metric)
	if err != nil {
		return err
	}

	valueHandler, err := valueHandlerAbstractFactory.GetValueHandler(metric)
	if err != nil {
		return err
	}

	metric = valueHandler.GetMutatedValueMetric(metric, metricFromStorage, isSet)
	err = repos.UpdateMetric(ctx, metric)

	var errExport error
	if config.Configuration.StoreInterval == 0 {
		errExport = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	}

	return multierr.Combine(err, errExport)
}

func (metricManager *MetricManager) ImportFrom(ctx context.Context, permanentStorage *interfaces.PermanentStorage) error {
	toImport, err := (*permanentStorage).Import()
	if err != nil {
		return err
	}

	repos := config.Configuration.Repository

	for _, metric := range toImport {
		if err = repos.UpdateMetric(ctx, metric); err != nil {
			return err
		}
	}

	return nil
}

func (metricManager *MetricManager) ExportTo(ctx context.Context, permanentStorage *interfaces.PermanentStorage) error {
	toExport, err := config.Configuration.Repository.GetAllMetrics(ctx)
	if err != nil {
		return err
	}

	if err = (*permanentStorage).Export(toExport); err != nil {
		return err
	}

	return nil
}

func (metricManager *MetricManager) Init() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategy)
	metricManager.client = client.NewClient()

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}
