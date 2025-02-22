package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type RandomValueStrategy struct{}

func (strategy *RandomValueStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	var value = metricData.RandomValue
	metric.MType = dictionary.GaugeMetricType
	metric.MName = dictionary.RandomValueMetricName
	metric.Value = &value

	return nil
}
