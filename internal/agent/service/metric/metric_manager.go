package metric

import (
	"context"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/buisness_logic/parser"
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

// MetricManager todo: Перенести код из client/main.go сюда.
type MetricManager struct {
	metricListToCollect []string
	dataCollector       interfaces.MetricDataCollectorInterface
}

var (
	ErrNoMetricDataCollector    = errors.New("no metric data collector found")
	ErrNotInitializedResultList = errors.New("not initialized metric list")
	ErrEmptyMetricNameList      = errors.New("empty metric name list")

	err                 error
	strategies          map[string]interfaces.ParsingStrategyInterface
	metric              dto.Metrics
	metricList          []dto.Metrics
	collectedMetricData agentDto.CollectedMetricData
)

func (m *MetricManager) Collect(ctx context.Context) (*[]dto.Metrics, error) {
	if nil == m.dataCollector {
		return nil, ErrNoMetricDataCollector
	}
	if nil == metricList {
		return nil, ErrNotInitializedResultList
	}

	collectedMetricData = agentDto.CollectedMetricData{}
	err = m.dataCollector.CollectMetricData(ctx, &collectedMetricData)
	if nil != err {
		return nil, err
	}

	metricList = metricList[:0]
	for _, metricName := range m.metricListToCollect {
		metric = dto.Metrics{}
		err = strategies[metricName].ParseMetric(&metric, &collectedMetricData)
		if nil != err {
			return nil, err
		}
		metricList = append(metricList, metric)
	}

	return &metricList, nil
}

func (m *MetricManager) initStrategies() error {
	strategies = make(map[string]interfaces.ParsingStrategyInterface)

	var err error
	for _, metricName := range m.metricListToCollect {
		strategies[metricName], err = parser.GetStrategy(metricName)
		if nil != err {
			return err
		}
	}

	return nil
}

func NewInstance(dataCollector interfaces.MetricDataCollectorInterface, metricNameList []string) (MetricManager, error) {
	if len(metricNameList) == 0 {
		return MetricManager{}, ErrEmptyMetricNameList
	}

	metricManager := MetricManager{
		metricListToCollect: metricNameList,
		dataCollector:       dataCollector,
	}

	metricList = make([]dto.Metrics, 0, len(metricList))

	err = metricManager.initStrategies()
	if nil != err {
		return MetricManager{}, err
	}

	return metricManager, nil
}
