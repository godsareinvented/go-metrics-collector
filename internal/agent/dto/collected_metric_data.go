package dto

import (
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
)

type CollectedMetricData struct {
	MemStats           runtime.MemStats
	VirtualMemoryStats mem.VirtualMemoryStat
	CPUPercentList     []float64
	PollCount          int64
	RandomValue        float64
}
