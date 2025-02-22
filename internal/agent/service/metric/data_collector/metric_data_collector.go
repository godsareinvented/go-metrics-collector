package data_collector

import (
	"github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"math/rand/v2"
	"runtime"
)

type MetricDataCollector struct {
	memStats  runtime.MemStats
	pollCount int64
}

func (metricCollector *MetricDataCollector) CollectMetricData(metricData *dto.CollectedMetricData) {
	metricCollector.pollCount += 1

	runtime.ReadMemStats(&metricCollector.memStats)

	metricData.PollCount = metricCollector.pollCount
	metricData.MemStats = metricCollector.memStats
	metricData.RandomValue = rand.Float64()
}

func New() *MetricDataCollector {
	return &MetricDataCollector{
		pollCount: 1,
	}
}
