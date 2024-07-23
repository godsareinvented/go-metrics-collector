package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type NumForcedGCStrategy[Num constraint.Numeric] struct{}

func (strategy *NumForcedGCStrategy[Num]) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num] {
	return dto.Metric[Num]{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: Num(metricData.MemStats.NumForcedGC), // float64
	}
}
