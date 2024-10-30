package manager

import (
	parserAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/value_handler/abstract_factory"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

// MetricManager todo: Стоит ли прорядить методы менеджера метрик?
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

func (metricManager *MetricManager) UpdateValue(metricDTO dto.Metrics) {
	repos := config.Configuration.Repository

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metricDTO, repos)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO)

	_, err := repos.UpdateMetric(metricDTO)
	if nil != err {
		panic("Error updating metric: " + err.Error())
	}
}

func (metricManager *MetricManager) GetByName(metric dto.Metrics) (dto.Metrics, bool) {
	repos := config.Configuration.Repository

	// todo: Вызов функции работает без разыменования?
	metricFromStorage, isSet, err := repos.GetMetricByName(metric)
	if isSet {
		return metricFromStorage, true
	}
	if nil != err {
		panic("Error getting metric: " + err.Error())
	}
	return metric, false
}

func (metricManager *MetricManager) GetByID(metric dto.Metrics) (dto.Metrics, bool) {
	repos := config.Configuration.Repository

	metricFromStorage, isSet, err := repos.GetMetricByID(metric)
	if isSet {
		return metricFromStorage, true
	}
	if nil != err {
		panic("Error getting metric: " + err.Error())
	}
	return metric, false
}

func (metricManager *MetricManager) GetList() []dto.Metrics {
	repos := config.Configuration.Repository

	metrics, err := repos.GetAllMetrics()
	if nil != err {
		panic("Error getting metric list: " + err.Error())
	}
	return metrics
}

func (metricManager *MetricManager) Init() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategyInterface)

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}
