package data_collector

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type CounterMetricDataCollector struct {
	pollCount int64
}

func (metricCollector *CounterMetricDataCollector) CollectMetricData(metricDataDTO *dto.CollectedMetricData) {
	metricCollector.pollCount += 1
	metricDataDTO.PollCount = metricCollector.pollCount
}
