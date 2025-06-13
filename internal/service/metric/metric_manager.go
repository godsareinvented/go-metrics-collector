package metric

import (
	"context"
	parserAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/parser"
	valueHandlerAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/value_handler"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"go.uber.org/multierr"
	"maps"
	"slices"
	"time"
)

type MetricManager struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
	Client              interfaces.Client
	strategies          map[string]interfaces.ParsingStrategy
}

var (
	metricList []dto.Metrics
)

func (metricManager *MetricManager) CollectAndSend(ctx context.Context) {
	if metricManager.MetricList == nil {
		panic("metric list is empty")
	}
	if metricManager.Client == nil {
		panic("Client is required")
	}

	go metricManager.collect(ctx)
	go metricManager.send(ctx)
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

			metrics = make([]dto.Metrics, 0, len(metricManager.MetricList))
			for _, metricName := range metricManager.MetricList {
				strategy := metricManager.strategies[metricName]
				m := strategy.GetMetric(metricName, metricCollectedData)
				metrics = append(metrics, m)
			}
			metricList = metrics

			time.Sleep(config.Configuration.PollInterval * time.Second)
		}
	}
}

func (metricManager *MetricManager) send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_ = metricManager.Client.SendBatch(metricList)

			time.Sleep(config.Configuration.ReportInterval * time.Second)
		}
	}
}

func (metricManager *MetricManager) UpdateMetric(ctx context.Context, metric dto.Metrics) error {
	repos := config.Configuration.Repository

	m, err := prepareMetric(ctx, *repos, metric)
	if err != nil {
		return err
	}
	err = repos.UpdateMetric(ctx, m)

	var errExport error
	if config.Configuration.StoreInterval == 0 {
		errExport = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	}

	return multierr.Combine(err, errExport)
}

func (metricManager *MetricManager) UpdateMetrics(ctx context.Context, metrics []dto.Metrics) ([]dto.Metrics, error) {
	repos := config.Configuration.Repository

	metrics, err := combineMetricValues(metrics)
	if err != nil {
		return nil, err
	}

	var resMetrics []dto.Metrics
	for _, metric := range metrics {
		m, err := prepareMetric(ctx, *repos, metric)
		if err != nil {
			return nil, err
		}
		resMetrics = append(resMetrics, m)
	}

	err = repos.UpdateMetricBatch(ctx, resMetrics)

	var errExport error
	if config.Configuration.StoreInterval == 0 {
		errExport = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	}

	return resMetrics, multierr.Combine(err, errExport)
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

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}

func prepareMetric(ctx context.Context, repos repository.Repository, metric dto.Metrics) (dto.Metrics, error) {
	metricFromStorage, isSet, err := repos.GetMetric(ctx, metric)
	if err != nil {
		return dto.Metrics{}, err
	}

	valueHandler, err := valueHandlerAbstractFactory.GetValueHandler(metric)
	if err != nil {
		return dto.Metrics{}, err
	}

	metric = valueHandler.GetMutatedValueMetric(metric, metricFromStorage, isSet)

	return metric, nil
}

func combineMetricValues(metricList []dto.Metrics) ([]dto.Metrics, error) {
	valueHandlerMap := map[string]interfaces.ValueHandler{}
	metricMap := map[string]dto.Metrics{}
	isMetricSet := true

	for _, metric := range metricList {
		isMetricSet = true
		if _, ok := metricMap[metric.ID]; !ok {
			isMetricSet = false
			if _, ok = valueHandlerMap[metric.MType]; !ok {
				valueHandler, err := valueHandlerAbstractFactory.GetValueHandler(metric)
				if err != nil {
					return []dto.Metrics{}, err
				}
				valueHandlerMap[metric.MType] = valueHandler
			}
		}
		metricMap[metric.ID] = valueHandlerMap[metric.MType].GetMutatedValueMetric(metric, metricMap[metric.ID], isMetricSet)
	}

	return slices.Collect(maps.Values(metricMap)), nil
}
