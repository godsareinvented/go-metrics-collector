package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
)

type MetricDataCollectorInterface interface {
	CollectMetricData(*dto.CollectedMetricData)
}
