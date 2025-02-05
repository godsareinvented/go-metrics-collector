package metric

import (
	"github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"math/rand"
	"runtime"
)

type MetricDataCollector struct{}

func (metricCollector *MetricDataCollector) CollectMetricData(metricDataDTO *dto.CollectedMetricData) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var pollCount int64 = 1

	metricDataDTO.PollCount = pollCount
	metricDataDTO.MemStats = memStats
	metricDataDTO.RandomValue = rand.Float64()
}
