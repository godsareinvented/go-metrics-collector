package metric

import (
	"github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"math/rand"
	"runtime"
)

type MetricDataCollector struct {
	memStats  runtime.MemStats
	pollCount int64
}

func (metricCollector *MetricDataCollector) CollectMetricData(metricData *dto.CollectedMetricData) {
	runtime.ReadMemStats(&metricCollector.memStats)

	metricData.PollCount = metricCollector.pollCount
	metricData.MemStats = metricCollector.memStats
	metricData.RandomValue = rand.Float64()
}

func NewDataCollector() *MetricDataCollector {
	return &MetricDataCollector{
		pollCount: 1,
	}
}
