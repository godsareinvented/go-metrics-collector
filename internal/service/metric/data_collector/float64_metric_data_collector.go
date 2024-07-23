package metric_data_collector

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"math/rand"
	"runtime"
)

type Float64MetricDataCollector struct{}

func (metricCollector *Float64MetricDataCollector) CollectMetricData(metricDataDTO *dto.CollectedMetricData) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricDataDTO.MemStats = memStats
	metricDataDTO.RandomValue = rand.Float64()
}
