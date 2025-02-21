package metric

import (
	"github.com/godsareinvented/go-metrics-collector/internal/agent/buisness_logic/parser"
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
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
	var collectedMetricData agentDto.CollectedMetricData

	metricManager.DataCollector.CollectMetricData(&collectedMetricData)

	for _, metricName := range metricManager.MetricList {
		metric = metricManager.strategies[metricName].GetMetric(metricName, collectedMetricData)
		metricList = append(metricList, metric)
	}

	return metricList
}

func (metricManager *MetricManager) Init() {
	metricManager.initStrategyList()
}

func (metricManager *MetricManager) initStrategyList() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategyInterface)

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parser.GetStrategy(metricName)
	}
}
