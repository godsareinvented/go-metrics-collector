package metric_data_collector

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"math/rand"
	"runtime"
)

type MetricDataCollector struct{}

func (metricCollector *MetricDataCollector) CollectMetricData(metricDataDTO *dto.CollectedMetricData) {
	var pollCount int64 = 1
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricDataDTO.PollCount = pollCount
	metricDataDTO.MemStats = memStats
	metricDataDTO.RandomValue = rand.Float64()
}
