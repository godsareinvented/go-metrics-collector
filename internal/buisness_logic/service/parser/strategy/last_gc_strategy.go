package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type LastGCStrategy[Num constraint.Numeric] struct{}

func (strategy *LastGCStrategy[Num]) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num] {
	return dto.Metric[Num]{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: Num(metricData.MemStats.LastGC), // float64
	}
}
