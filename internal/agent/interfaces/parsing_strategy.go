package interfaces

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type ParsingStrategy interface {
	GetMetric(metricName string, metricData agentdto.CollectedMetricData) generaldto.Metrics
}
