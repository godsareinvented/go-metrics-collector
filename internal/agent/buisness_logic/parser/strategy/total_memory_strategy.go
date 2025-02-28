package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type TotalMemoryStrategy struct{}

func (strategy *TotalMemoryStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	if isVirtualMemoryStatsMemEmpty(&metricData.VirtualMemoryStats) {
		return ErrEmptyCollectedData
	}

	var value = float64(metricData.VirtualMemoryStats.Total)
	metric.MType = dictionary.GaugeMetricType
	metric.MName = dictionary.TotalMemoryMetricName
	metric.Value = &value

	return nil
}
