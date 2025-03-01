package interfaces

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type ParsingStrategy interface {
	FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData)
}
