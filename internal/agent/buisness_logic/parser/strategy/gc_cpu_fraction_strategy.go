package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type GCCPUFractionStrategy struct{}

func (strategy *GCCPUFractionStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	if isMemStatEmpty(&metricData.MemStats) {
		return ErrEmptyCollectedData
	}

	var value = metricData.MemStats.GCCPUFraction
	metric.MType = dictionary.GaugeMetricType
	metric.MName = dictionary.GCCPUFractionMetricName
	metric.Value = &value

	return nil
}
