package metric

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
	"sync"
)

type MetricDataCollector struct {
	wg                sync.WaitGroup
	virtualMemoryStat *mem.VirtualMemoryStat
	pollCount         int64
}

func (collector *MetricDataCollector) CollectMetricData(ctx context.Context, metricData *dto.CollectedMetricData) error {
	collector.wg.Add(2)

	var mcErr error
	go collector.CollectMemStats(metricData)
	go collector.CollectVirtualMemoryStatsAndCpuUtilization(ctx, metricData, &mcErr)
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

func (collector *MetricDataCollector) CollectVirtualMemoryStatsAndCpuUtilization(ctx context.Context, metricData *dto.CollectedMetricData, parentErr *error) {
	virtualMemoryStat, err := mem.VirtualMemory()
	if nil != err {
		*parentErr = err
		collector.wg.Done()
		return
	}

	metricData.CPUPercentList, err = cpu.PercentWithContext(ctx, 0, true)
	if nil != err {
		collector.wg.Done()
		return
	}

	metricData.VirtualMemoryStats = *virtualMemoryStat
	collector.wg.Done()
}

func NewDataCollector() *MetricDataCollector {
	return &MetricDataCollector{
		pollCount: 1,
	}
}
