package interfaces

import (
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type ValueHandler interface {
	UpdateMetricValue(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics
}
