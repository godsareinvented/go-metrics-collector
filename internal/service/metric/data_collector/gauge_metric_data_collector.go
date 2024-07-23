package data_collector

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"math/rand"
	"runtime"
)

type GaugeMetricDataCollector struct{}

func (metricCollector *GaugeMetricDataCollector) CollectMetricData(metricData *dto.CollectedMetricData) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricData.MemStats = memStats
	metricData.RandomValue = rand.Float64()
}
