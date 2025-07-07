package parser

import (
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/parser/strategy"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/parser/strategy/cpu"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/parser/strategy/mem"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/parser/strategy/virtual_mem"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"strconv"
	"strings"
)

var (
	ErrUnknownMetricName = errors.New("unknown metric name")

	strategyMap = map[string]interfaces.ParsingStrategy{
		dictionary.AllocMetricName:         &mem.AllocStrategy{},
		dictionary.BuckHashSysMetricName:   &mem.BuckHashSysStrategy{},
		dictionary.FreesMetricName:         &mem.FreesStrategy{},
		dictionary.GCCPUFractionMetricName: &mem.GCCPUFractionStrategy{},
		dictionary.GCSysMetricName:         &mem.GCSysStrategy{},
		dictionary.HeapAllocMetricName:     &mem.HeapAllocStrategy{},
		dictionary.HeapIdleMetricName:      &mem.HeapIdleStrategy{},
		dictionary.HeapInuseMetricName:     &mem.HeapInuseStrategy{},
		dictionary.HeapObjectsMetricName:   &mem.HeapObjectsStrategy{},
		dictionary.HeapReleasedMetricName:  &mem.HeapReleasedStrategy{},
		dictionary.HeapSysMetricName:       &mem.HeapSysStrategy{},
		dictionary.LastGCMetricName:        &mem.LastGCStrategy{},
		dictionary.LookupsMetricName:       &mem.LookupsStrategy{},
		dictionary.MCacheInuseMetricName:   &mem.MCacheInuseStrategy{},
		dictionary.MCacheSysMetricName:     &mem.MCacheSysStrategy{},
		dictionary.MSpanInuseMetricName:    &mem.MSpanInuseStrategy{},
		dictionary.MSpanSysMetricName:      &mem.MSpanSysStrategy{},
		dictionary.MallocsMetricName:       &mem.MallocsStrategy{},
		dictionary.NextGCMetricName:        &mem.NextGCStrategy{},
		dictionary.NumForcedGCMetricName:   &mem.NumForcedGCStrategy{},
		dictionary.NumGCMetricName:         &mem.NumGCStrategy{},
		dictionary.OtherSysMetricName:      &mem.OtherSysStrategy{},
		dictionary.PauseTotalNsMetricName:  &mem.PauseTotalNsStrategy{},
		dictionary.StackInuseMetricName:    &mem.StackInuseStrategy{},
		dictionary.StackSysMetricName:      &mem.StackSysStrategy{},
		dictionary.SysMetricName:           &mem.SysStrategy{},
		dictionary.TotalAllocMetricName:    &mem.TotalAllocStrategy{},
		dictionary.PollCountMetricName:     &strategy.PollCountStrategy{},
		dictionary.RandomValueMetricName:   &strategy.RandomValueStrategy{},
		dictionary.TotalMemoryMetricName:   &virtual_mem.TotalMemoryStrategy{},
		dictionary.FreeMemoryMetricName:    &virtual_mem.FreeMemoryStrategy{},
	}
)

func Strategy(metricName string) (interfaces.ParsingStrategy, error) {
	if parsingStrategy, ok := strategyFromMap(metricName); ok {
		return parsingStrategy, nil
	}

	if parsingStrategy, ok := parsedCpuUtilizationStrategy(metricName); ok {
		return parsingStrategy, nil
	}

	return nil, ErrUnknownMetricName
}

func strategyFromMap(metricName string) (interfaces.ParsingStrategy, bool) {
	if _, ok := strategyMap[metricName]; ok {
		return strategyMap[metricName], true
	}
	return nil, false
}

// CPUutilizaition0, CPUutilizaition1,..
// ref: dictionary.CPUutilizationMetricName
func parsedCpuUtilizationStrategy(metricName string) (interfaces.ParsingStrategy, bool) {
	logicalCpuNumber, ok := getLogicalCpuNumberFromMetricName(metricName)
	if !ok {
		return nil, false
	}
	return cpu.NewCpuUtilizationStrategy(logicalCpuNumber, metricName), true
}

func getLogicalCpuNumberFromMetricName(metricName string) (uint, bool) {
	subStrings := strings.Split(metricName, dictionary.CpuUtilizationMetricName)
	if len(subStrings) != 2 {
		return 0, false
	}
	logicalCpuNumber, err := strconv.ParseUint(subStrings[1], 10, 8)
	if nil != err {
		return 0, false
	}
	return uint(logicalCpuNumber), true
}
