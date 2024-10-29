package metric_data_collector

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"math/rand"
	"runtime"
)

type MetricDataCollector struct{}

var pollCount int64 = 1

func (metricCollector *MetricDataCollector) CollectMetricData(metricDataDTO *dto.CollectedMetricData) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	pollCount += 1

	metricDataDTO.PollCount = pollCount
	metricDataDTO.MemStats = memStats
	metricDataDTO.RandomValue = rand.Float64()
}
