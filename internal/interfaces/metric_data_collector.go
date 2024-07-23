package interfaces

import "github.com/godsareinvented/go-metrics-collector/internal/dto"

type MetricDataCollector interface {
	CollectMetricData(*dto.CollectedMetricData)
}
