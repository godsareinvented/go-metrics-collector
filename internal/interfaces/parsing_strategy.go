package interfaces

import "github.com/oldhanasong/go-metrics-collector/internal/dto"

type ParsingStrategy interface {
	GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics
}
