package data_collector

import (
	"github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"runtime"
	"testing"
)

const epsilon float64 = 0.1    // Допустимая ошибка - 10%
const pollCountDelta int64 = 1 // Допустимая разница для PollCount

func TestGetMetricData(t *testing.T) {
	var runtimeData dto.CollectedMetricData
	var collectedData dto.CollectedMetricData
	var secondTimeCollectedData dto.CollectedMetricData

	collector := MetricDataCollector{}

	runtime.ReadMemStats(&runtimeData.MemStats)

	err := collector.Collect(&collectedData)
	require.NoError(t, err)

	err = collector.Collect(&secondTimeCollectedData)
	require.NoError(t, err)

	testMemoryStatsCollecting(t, collectedData.MemStats, runtimeData.MemStats)
	testPollCount(t, collectedData.PollCount, secondTimeCollectedData.PollCount)
}

func testMemoryStatsCollecting(t *testing.T, collected, target runtime.MemStats) {
	t.Run("mem stats collecting", func(t *testing.T) {
		tests := []struct {
			metricName string
			collected  float64
			target     float64
		}{
			{metricName: dictionary.AllocMetricName, collected: float64(collected.Alloc), target: float64(target.Alloc)},
			{metricName: dictionary.BuckHashSysMetricName, collected: float64(collected.BuckHashSys), target: float64(target.BuckHashSys)},
			{metricName: dictionary.FreesMetricName, collected: float64(collected.Frees), target: float64(target.Frees)},
			{metricName: dictionary.GCCPUFractionMetricName, collected: collected.GCCPUFraction, target: target.GCCPUFraction},
			{metricName: dictionary.GCSysMetricName, collected: float64(collected.GCSys), target: float64(target.GCSys)},
			{metricName: dictionary.HeapAllocMetricName, collected: float64(collected.HeapAlloc), target: float64(target.HeapAlloc)},
			{metricName: dictionary.HeapIdleMetricName, collected: float64(collected.HeapIdle), target: float64(target.HeapIdle)},
			{metricName: dictionary.HeapInuseMetricName, collected: float64(collected.HeapInuse), target: float64(target.HeapInuse)},
			{metricName: dictionary.HeapObjectsMetricName, collected: float64(collected.HeapObjects), target: float64(target.HeapObjects)},
			{metricName: dictionary.HeapReleasedMetricName, collected: float64(collected.HeapReleased), target: float64(target.HeapReleased)},
			{metricName: dictionary.HeapSysMetricName, collected: float64(collected.HeapSys), target: float64(target.HeapSys)},
			{metricName: dictionary.LastGCMetricName, collected: float64(collected.LastGC), target: float64(target.LastGC)},
			{metricName: dictionary.LookupsMetricName, collected: float64(collected.Lookups), target: float64(target.Lookups)},
			{metricName: dictionary.MCacheInuseMetricName, collected: float64(collected.MCacheInuse), target: float64(target.MCacheInuse)},
			{metricName: dictionary.MCacheSysMetricName, collected: float64(collected.MCacheSys), target: float64(target.MCacheSys)},
			{metricName: dictionary.MSpanInuseMetricName, collected: float64(collected.MSpanInuse), target: float64(target.MSpanInuse)},
			{metricName: dictionary.MSpanSysMetricName, collected: float64(collected.MSpanSys), target: float64(target.MSpanSys)},
			{metricName: dictionary.MallocsMetricName, collected: float64(collected.Mallocs), target: float64(target.Mallocs)},
			{metricName: dictionary.NextGCMetricName, collected: float64(collected.NextGC), target: float64(target.NextGC)},
			{metricName: dictionary.NumForcedGCMetricName, collected: float64(collected.NumForcedGC), target: float64(target.NumForcedGC)},
			{metricName: dictionary.NumGCMetricName, collected: float64(collected.NumGC), target: float64(target.NumGC)},
			{metricName: dictionary.OtherSysMetricName, collected: float64(collected.OtherSys), target: float64(target.OtherSys)},
			{metricName: dictionary.PauseTotalNsMetricName, collected: float64(collected.PauseTotalNs), target: float64(target.PauseTotalNs)},
			{metricName: dictionary.StackInuseMetricName, collected: float64(collected.StackInuse), target: float64(target.StackInuse)},
			{metricName: dictionary.StackSysMetricName, collected: float64(collected.StackSys), target: float64(target.StackSys)},
			{metricName: dictionary.SysMetricName, collected: float64(collected.Sys), target: float64(target.Sys)},
			{metricName: dictionary.TotalAllocMetricName, collected: float64(collected.TotalAlloc), target: float64(target.TotalAlloc)},
		}

		for _, test := range tests {
			if test.target == test.collected {
				continue
			}
			change := math.Abs((test.target - test.collected) / test.collected)
			assert.Truef(
				t,
				change <= epsilon,
				"The change of %d%% is outside the range of ±%d%% (old: %.2f, new: %.2f)",
				int(change*100),
				int(epsilon*100),
				test.collected,
				test.target,
			)
		}
	})
}

func testPollCount(t *testing.T, firstTimeValue, secondTimeValue int64) {
	t.Run("increment PollCount", func(t *testing.T) {
		delta := secondTimeValue - firstTimeValue
		require.LessOrEqualf(
			t,
			delta,
			pollCountDelta,
			"PollCount mismatch: first time got %d, second time got %d (needed value for second time is %d)",
			firstTimeValue,
			secondTimeValue,
			firstTimeValue+pollCountDelta,
		)
	})
}
