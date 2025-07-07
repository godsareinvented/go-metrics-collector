package data_collector

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand/v2"
	"runtime"
	"sync"
)

type MetricDataCollector struct {
	wg                sync.WaitGroup
	virtualMemoryStat *mem.VirtualMemoryStat
	pollCount         int64
}

func (collector *MetricDataCollector) Collect(metricData *dto.CollectedMetricData) error {
	collector.wg.Add(2)

	var mcErr error
	go collector.CollectMemStats(metricData)
	go collector.CollectVirtualMemoryStatsAndCpuUtilization(metricData, &mcErr)
	metricData.PollCount = collector.pollCount
	metricData.RandomValue = rand.Float64()

	collector.wg.Wait()

	if nil != mcErr {
		return mcErr
	}
	return nil
}

func (collector *MetricDataCollector) CollectMemStats(metricData *dto.CollectedMetricData) {
	runtime.ReadMemStats(&metricData.MemStats)
	collector.wg.Done()
}

func (collector *MetricDataCollector) CollectVirtualMemoryStatsAndCpuUtilization(metricData *dto.CollectedMetricData, parentErr *error) {
	defer collector.wg.Done()

	virtualMemoryStat, err := mem.VirtualMemory()
	if nil != err {
		*parentErr = err
		return
	}
	metricData.VirtualMemoryStats = *virtualMemoryStat

	if metricData.CPUPercentList, err = cpu.PercentWithContext(context.Background(), 0, true); nil != err {
		*parentErr = err
		return
	}
}

func New() *MetricDataCollector {
	return &MetricDataCollector{
		pollCount: 1,
	}
}
