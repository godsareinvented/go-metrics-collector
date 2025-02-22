package interfaces

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type ParsingStrategyInterface interface {
	ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error
}
