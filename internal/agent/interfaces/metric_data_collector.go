package interfaces

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
)

type MetricDataCollectorInterface interface {
	CollectMetricData(ctx context.Context, metricData *dto.CollectedMetricData) error
}
