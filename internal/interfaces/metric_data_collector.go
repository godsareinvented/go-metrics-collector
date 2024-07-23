package interfaces

import "github.com/oldhanasong/go-metrics-collector/internal/dto"

type MetricDataCollector interface {
	CollectMetricData(*dto.CollectedMetricData)
}
