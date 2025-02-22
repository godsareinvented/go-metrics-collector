package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type TotalAllockStrategy struct{}

func (strategy *TotalAllockStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	if isMemStatEmpty(&metricData.MemStats) {
		return ErrEmptyCollectedData
	}

	var value = float64(metricData.MemStats.TotalAlloc)
	metric.MType = dictionary.GaugeMetricType
	metric.MName = dictionary.TotalAllocMetricName
	metric.Value = &value

	return nil
}
