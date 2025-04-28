package metric

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/server/buisness_logic/value_handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/server/repository"
)

type MetricManager struct{}

func (metricManager *MetricManager) UpdateMetric(ctx context.Context, metric dto.Metrics) error {
	repos := config.Configuration.Repository

	metric, err := metricManager.getPreparedMetric(ctx, *repos, &metric)
	if nil != err {
		return err
	}

	_, err = repos.UpdateMetric(ctx, metric)
	if nil != err {
		return err
	}

	if 0 == config.Configuration.StoreInterval {
		err = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
		if nil != err {
			return err
		}
	}

	return nil
}

func (metricManager *MetricManager) UpdateMetrics(ctx context.Context, metrics []dto.Metrics) error {
	repos := config.Configuration.Repository

	var metric dto.Metrics
	var resultingMetrics []dto.Metrics
	var err error
	for _, metric = range metrics {
		metric, err = metricManager.getPreparedMetric(ctx, *repos, &metric)
		if nil != err {
			return err
		}
		resultingMetrics = append(resultingMetrics, metric)
	}

	err = repos.UpdateMetricBatch(ctx, resultingMetrics)
	if nil != err {
		return err
	}

	if 0 == config.Configuration.StoreInterval {
		err = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
		if nil != err {
			return err
		}
	}

	return nil
}

func (metricManager *MetricManager) ImportFrom(ctx context.Context, permanentStorage *interfaces.PermanentStorage) error {
	metricList, err := (*permanentStorage).Import()
	if nil != err {
		return err
	}
	err = config.Configuration.Repository.UpdateMetricBatch(ctx, metricList)
	if nil != err {
		return err
	}

	return nil
}

func (metricManager *MetricManager) ExportTo(ctx context.Context, permanentStorage *interfaces.PermanentStorage) error {
	metricList, err := config.Configuration.Repository.GetAllMetrics(ctx)
	if nil != err {
		return err
	}

	err = (*permanentStorage).Export(metricList)
	if err != nil {
		return err
	}

	return nil
}

// getPreparedMetric todo: Необходимо убрать лишние запросы на каждую отдельную метрику и заменить их единым запросом
func (metricManager *MetricManager) getPreparedMetric(ctx context.Context, repos repository.Repository, metric *dto.Metrics) (dto.Metrics, error) {
	metricFromStorage, isSet, err := repos.GetMetric(ctx, *metric)
	if nil != err {
		return dto.Metrics{}, err
	}

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(*metric)
	mutatedValueMetric := valueHandler.GetMutatedValueMetric(*metric, metricFromStorage, isSet)

	if "" != metricFromStorage.ID {
		mutatedValueMetric.ID = metricFromStorage.ID
	}

	return mutatedValueMetric, nil
}
