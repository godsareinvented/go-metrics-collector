package data_collector

import (
	"github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"math/rand/v2"
	"runtime"
)

type MetricDataCollector struct {
	pollCount int64
}

func (metricCollector *MetricDataCollector) CollectMetricData(metricData *dto.CollectedMetricData) {
	metricCollector.pollCount += 1

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricData.PollCount = metricCollector.pollCount
	metricData.MemStats = memStats
	metricData.RandomValue = rand.Float64()
}
