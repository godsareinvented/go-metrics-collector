package interfaces

import (
	dto2 "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type ParsingStrategyInterface interface {
	GetMetric(metricName string, metricData dto2.CollectedMetricData) dto.Metrics
}
