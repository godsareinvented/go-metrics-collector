package strategy

import (
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
)

var (
	ErrEmptyCollectedData = errors.New("empty collected data")

	emptyMemStatsStruct           = runtime.MemStats{}
	emptyVirtualMemoryStatsStruct = mem.VirtualMemoryStat{}
)

func isMemStatEmpty(memStats *runtime.MemStats) bool {
	return *memStats == emptyMemStatsStruct
}

func isVirtualMemoryStatsMemEmpty(virtualMemoryStats *mem.VirtualMemoryStat) bool {
	return *virtualMemoryStats == emptyVirtualMemoryStatsStruct
}

func isLogicalCpuPercentsValid(logicalCpuPercents *[]float64) bool {
	return len(*logicalCpuPercents) == config.Configuration.LogicalCpuCount
}
