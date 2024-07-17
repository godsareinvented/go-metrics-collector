package metric_data_collector

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"math/rand"
	"runtime"
)

type MetricDataCollector struct{}

func (metricCollector *MetricDataCollector) GetMetricData() dto.CollectedMetricData {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var pollCount int64 = 1

	return dto.CollectedMetricData{MemStats: memStats, PollCount: pollCount, RandomValue: rand.Float64()}
}
