package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type MCacheSysStrategy struct{}

func (strategy *MCacheSysStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	if isMemStatEmpty(&metricData.MemStats) {
		return ErrEmptyCollectedData
	}

	var value = float64(metricData.MemStats.MCacheSys)
	metric.MType = dictionary.GaugeMetricType
	metric.MName = dictionary.MCacheSysMetricName
	metric.Value = &value

	return nil
}
