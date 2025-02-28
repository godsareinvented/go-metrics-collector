package parser

import (
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/buisness_logic/parser/strategy"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	"strconv"
	"strings"
)

var (
	ErrUnknownMetricName = errors.New("unknown metric name")

	strategyMap = map[string]interfaces.ParsingStrategy{
		dictionary.AllocMetricName:         &strategy.AllocStrategy{},
		dictionary.BuckHashSysMetricName:   &strategy.BuckHashSysStrategy{},
		dictionary.FreesMetricName:         &strategy.FreesStrategy{},
		dictionary.GCCPUFractionMetricName: &strategy.GCCPUFractionStrategy{},
		dictionary.GCSysMetricName:         &strategy.GCSysStrategy{},
		dictionary.HeapAllocMetricName:     &strategy.HeapAllocStrategy{},
		dictionary.HeapIdleMetricName:      &strategy.HeapIdleStrategy{},
		dictionary.HeapInuseMetricName:     &strategy.HeapInuseStrategy{},
		dictionary.HeapObjectsMetricName:   &strategy.HeapObjectsStrategy{},
		dictionary.HeapReleasedMetricName:  &strategy.HeapReleasedStrategy{},
		dictionary.HeapSysMetricName:       &strategy.HeapSysStrategy{},
		dictionary.LastGCMetricName:        &strategy.LastGCStrategy{},
		dictionary.LookupsMetricName:       &strategy.LookupsStrategy{},
		dictionary.MCacheInuseMetricName:   &strategy.MCacheInuseStrategy{},
		dictionary.MCacheSysMetricName:     &strategy.MCacheSysStrategy{},
		dictionary.MSpanInuseMetricName:    &strategy.MSpanInuseStrategy{},
		dictionary.MSpanSysMetricName:      &strategy.MSpanSysStrategy{},
		dictionary.MallocsMetricName:       &strategy.MallocsStrategy{},
		dictionary.NextGCMetricName:        &strategy.NextGCStrategy{},
		dictionary.NumForcedGCMetricName:   &strategy.NumForcedGCStrategy{},
		dictionary.NumGCMetricName:         &strategy.NumGCStrategy{},
		dictionary.OtherSysMetricName:      &strategy.OtherSysStrategy{},
		dictionary.PauseTotalNsMetricName:  &strategy.PauseTotalNsStrategy{},
		dictionary.StackInuseMetricName:    &strategy.StackInuseStrategy{},
		dictionary.StackSysMetricName:      &strategy.StackSysStrategy{},
		dictionary.SysMetricName:           &strategy.SysStrategy{},
		dictionary.TotalAllocMetricName:    &strategy.TotalAllocStrategy{},
		dictionary.PollCountMetricName:     &strategy.PollCountStrategy{},
		dictionary.RandomValueMetricName:   &strategy.RandomValueStrategy{},
		dictionary.TotalMemoryMetricName:   &strategy.TotalMemoryStrategy{},
		dictionary.FreeMemoryMetricName:    &strategy.FreeMemoryStrategy{},
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
	return strategy.NewCpuUtilizationStrategy(logicalCpuNumber, metricName), true
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
