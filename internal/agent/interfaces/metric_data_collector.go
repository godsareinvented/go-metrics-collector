package interfaces

import (
	"github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
)

type MetricDataCollector interface {
	CollectMetricData(*dto.CollectedMetricData)
}
