package manager

import (
	"context"
	parserAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/value_handler/abstract_factory"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type MetricManager struct {
	MetricList    []string
	DataCollector interfaces.MetricDataCollectorInterface
	strategies    map[string]interfaces.ParsingStrategyInterface
}

func (metricManager *MetricManager) Collect() []dto.Metrics {
	if nil == metricManager.DataCollector {
		panic("nil DataCollector")
	}

	var metric dto.Metrics
	var metricList []dto.Metrics
	var collectedMetricData dto.CollectedMetricData

	metricManager.DataCollector.CollectMetricData(&collectedMetricData)

	for _, metricName := range metricManager.MetricList {
		metric = metricManager.strategies[metricName].GetMetric(metricName, collectedMetricData)
		metricList = append(metricList, metric)
	}

	return metricList
}

func (metricManager *MetricManager) UpdateMetric(ctx context.Context, metric dto.Metrics) {
	repos := config.Configuration.Repository

	if _, err := repos.UpdateMetric(ctx, metricManager.getPreparedMetric(ctx, *repos, &metric)); nil != err {
		// todo: Надо пересмотреть выплёвывание ошибок.
		panic("Error updating metric: " + err.Error())
	}

	if 0 == config.Configuration.StoreInterval {
		_ = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	}
}

func (metricManager *MetricManager) UpdateMetrics(ctx context.Context, metrics []dto.Metrics) {
	repos := config.Configuration.Repository

	var resultingMetrics []dto.Metrics
	for _, metric := range metrics {
		resultingMetrics = append(resultingMetrics, metricManager.getPreparedMetric(ctx, *repos, &metric))
	}

	if err := repos.UpdateMetricBatch(ctx, resultingMetrics); nil != err {
		panic("Error updating metric: " + err.Error())
	}

	if 0 == config.Configuration.StoreInterval {
		_ = metricManager.ExportTo(ctx, config.Configuration.PermanentStorage)
	}
}

func (metricManager *MetricManager) ImportFrom(ctx context.Context, permanentStorage *interfaces.PermanentStorage) error {
	metricList, err := (*permanentStorage).Import()
	if nil != err {
		return err
	}

	for _, metric := range metricList {
		metricManager.UpdateMetric(ctx, metric)
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

func (metricManager *MetricManager) Init() {
	metricManager.initStrategyList()
}

func (metricManager *MetricManager) initStrategyList() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategyInterface)

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}

func (metricManager *MetricManager) getPreparedMetric(ctx context.Context, repos repository.Repository, metric *dto.Metrics) dto.Metrics {
	metricFromStorage, isSet, _ := repos.GetMetric(ctx, *metric)

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(*metric)
	mutatedValueMetric := valueHandler.GetMutatedValueMetric(*metric, metricFromStorage, isSet)

	if "" != metricFromStorage.ID {
		mutatedValueMetric.ID = metricFromStorage.ID
	}

	return mutatedValueMetric
}
