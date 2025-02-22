package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type BuckHashSysStrategy struct{}

func (strategy *BuckHashSysStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	if isMemStatEmpty(&metricData.MemStats) {
		return ErrEmptyCollectedData
	}

	var value = float64(metricData.MemStats.BuckHashSys)
	metric.MType = dictionary.GaugeMetricType
	metric.MName = dictionary.BuckHashSysMetricName
	metric.Value = &value

	return nil
}
