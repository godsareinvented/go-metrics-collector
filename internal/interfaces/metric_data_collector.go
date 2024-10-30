package interfaces

import "github.com/godsareinvented/go-metrics-collector/internal/dto"

type MetricDataCollectorInterface interface {
	CollectMetricData(*dto.CollectedMetricData)
}
