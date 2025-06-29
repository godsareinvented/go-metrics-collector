package metric

import (
	"context"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	valueHandlerFactory "github.com/oldhanasong/go-metrics-collector/internal/server/buisness_logic/value_handler"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
	"go.uber.org/multierr"
	"maps"
	"slices"
)

type MetricManager struct{}

func (metricManager *MetricManager) UpdateMetric(ctx context.Context, metric generaldto.Metrics) error {
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

func (metricManager *MetricManager) UpdateMetrics(ctx context.Context, metrics []generaldto.Metrics) ([]generaldto.Metrics, error) {
	repos := config.Configuration.Repository

	metrics, err := combineMetricValues(metrics)
	if err != nil {
		return nil, err
	}

	var resMetrics []generaldto.Metrics
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

	err = config.Configuration.Repository.UpdateMetricBatch(ctx, toImport)
	if err != nil {
		return err
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

func prepareMetric(ctx context.Context, repos repository.Repository, metric generaldto.Metrics) (generaldto.Metrics, error) {
	metricFromStorage, isSet, err := repos.GetMetric(ctx, metric)
	if err != nil {
		return generaldto.Metrics{}, err
	}

	valueHandler, err := valueHandlerFactory.GetValueHandler(metric)
	if err != nil {
		return generaldto.Metrics{}, err
	}

	metric = valueHandler.UpdateMetricValue(metric, metricFromStorage, isSet)

	return metric, nil
}

func combineMetricValues(metricList []generaldto.Metrics) ([]generaldto.Metrics, error) {
	valueHandlerMap := map[string]interfaces.ValueHandler{}
	metricMap := map[string]generaldto.Metrics{}
	isMetricSet := true

	for _, metric := range metricList {
		isMetricSet = true
		if _, ok := metricMap[metric.ID]; !ok {
			isMetricSet = false
			if _, ok = valueHandlerMap[metric.MType]; !ok {
				valueHandler, err := valueHandlerFactory.GetValueHandler(metric)
				if err != nil {
					return []generaldto.Metrics{}, err
				}
				valueHandlerMap[metric.MType] = valueHandler
			}
		}
		metricMap[metric.ID] = valueHandlerMap[metric.MType].UpdateMetricValue(metric, metricMap[metric.ID], isMetricSet)
	}

	return slices.Collect(maps.Values(metricMap)), nil
}
