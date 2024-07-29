package manager

import (
	parserAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/value_handler/abstract_factory"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type MetricManager struct {
	MetricList    []string
	DataCollector interfaces.MetricDataCollector
	strategies    map[string]interfaces.ParsingStrategy
}

func (metricManager *MetricManager) Collect() []dto.Metric {
	if nil == metricManager.DataCollector {
		panic("nil DataCollector")
	}

	var metricDTO dto.Metric
	var metricDTOList []dto.Metric
	var collectedMetricData dto.CollectedMetricData

	metricManager.DataCollector.CollectMetricData(&collectedMetricData)

	for _, metricName := range metricManager.MetricList {
		metricDTO = metricManager.strategies[metricName].GetMetric(metricName, collectedMetricData)
		metricDTOList = append(metricDTOList, metricDTO)
	}

	return metricDTOList
}

func (metricManager *MetricManager) UpdateValue(metricDTO dto.Metric) {
	repos := config.Configuration.Repository

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metricDTO, repos)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO)

	repos.UpdateMetric(metricDTO)
}

func (metricManager *MetricManager) Get(metricDTO dto.Metric) (dto.Metric, bool) {
	repos := config.Configuration.Repository

	// todo: Вызов функции работает без разыменования?
	metricDTOFromDb, isSet := repos.GetMetric(metricDTO)
	if isSet {
		return metricDTOFromDb, true
	}
	return metricDTO, false
}

func (metricManager *MetricManager) GetList() []dto.Metric {
	repos := config.Configuration.Repository

	return repos.GetAllMetrics()
}

func (metricManager *MetricManager) Init() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategy)

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}
