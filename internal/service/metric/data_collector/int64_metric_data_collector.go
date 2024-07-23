package metric_data_collector

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type Int64MetricDataCollector struct{}

func (metricCollector *Int64MetricDataCollector) CollectMetricData(metricDataDTO *dto.CollectedMetricData) {
	var pollCount int64 = 1
	metricDataDTO.PollCount = pollCount
}
