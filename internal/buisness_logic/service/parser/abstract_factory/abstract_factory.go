package abstract_factory

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/parser/strategy"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

func GetStrategy(metricName string) interfaces.ParsingStrategy {
	strategyMap := map[string]interfaces.ParsingStrategy{
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
		dictionary.TotalAllocMetricName:    &strategy.TotalAllockStrategy{},
		dictionary.PollCountMetricName:     &strategy.PollCountStrategy{},
		dictionary.RandomValueMetricName:   &strategy.RandomValueStrategy{},
	}

	if _, ok := strategyMap[metricName]; !ok {
		panic("unknown metric type")
	}

	return strategyMap[metricName]
}
