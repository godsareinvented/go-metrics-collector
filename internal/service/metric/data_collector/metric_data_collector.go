package data_collector

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"math/rand"
	"runtime"
)

type MetricDataCollector struct {
	pollCount int64
}

func (metricCollector *MetricDataCollector) GetMetricData() dto.CollectedMetricData {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricCollector.pollCount += 1

	return dto.CollectedMetricData{MemStats: memStats, PollCount: metricCollector.pollCount, RandomValue: rand.Float64()}
}
