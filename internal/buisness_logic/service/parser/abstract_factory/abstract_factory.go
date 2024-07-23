package abstract_factory

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/parser/strategy"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

func GetStrategy[Num constraint.Numeric](metricName string) interfaces.ParsingStrategy[Num] {
	strategyMap := map[string]interfaces.ParsingStrategy[Num]{
		dictionary.AllocMetricName:         &strategy.AllocStrategy[Num]{},
		dictionary.BuckHashSysMetricName:   &strategy.BuckHashSysStrategy[Num]{},
		dictionary.FreesMetricName:         &strategy.FreesStrategy[Num]{},
		dictionary.GCCPUFractionMetricName: &strategy.GCCPUFractionStrategy[Num]{},
		dictionary.GCSysMetricName:         &strategy.GCSysStrategy[Num]{},
		dictionary.HeapAllocMetricName:     &strategy.HeapAllocStrategy[Num]{},
		dictionary.HeapIdleMetricName:      &strategy.HeapIdleStrategy[Num]{},
		dictionary.HeapInuseMetricName:     &strategy.HeapInuseStrategy[Num]{},
		dictionary.HeapObjectsMetricName:   &strategy.HeapObjectsStrategy[Num]{},
		dictionary.HeapReleasedMetricName:  &strategy.HeapReleasedStrategy[Num]{},
		dictionary.HeapSysMetricName:       &strategy.HeapSysStrategy[Num]{},
		dictionary.LastGCMetricName:        &strategy.LastGCStrategy[Num]{},
		dictionary.LookupsMetricName:       &strategy.LookupsStrategy[Num]{},
		dictionary.MCacheInuseMetricName:   &strategy.MCacheInuseStrategy[Num]{},
		dictionary.MCacheSysMetricName:     &strategy.MCacheSysStrategy[Num]{},
		dictionary.MSpanInuseMetricName:    &strategy.MSpanInuseStrategy[Num]{},
		dictionary.MSpanSysMetricName:      &strategy.MSpanSysStrategy[Num]{},
		dictionary.MallocsMetricName:       &strategy.MallocsStrategy[Num]{},
		dictionary.NextGCMetricName:        &strategy.NextGCStrategy[Num]{},
		dictionary.NumForcedGCMetricName:   &strategy.NumForcedGCStrategy[Num]{},
		dictionary.NumGCMetricName:         &strategy.NumGCStrategy[Num]{},
		dictionary.OtherSysMetricName:      &strategy.OtherSysStrategy[Num]{},
		dictionary.PauseTotalNsMetricName:  &strategy.PauseTotalNsStrategy[Num]{},
		dictionary.StackInuseMetricName:    &strategy.StackInuseStrategy[Num]{},
		dictionary.StackSysMetricName:      &strategy.StackSysStrategy[Num]{},
		dictionary.SysMetricName:           &strategy.SysStrategy[Num]{},
		dictionary.TotalAllocMetricName:    &strategy.TotalAllockStrategy[Num]{},
		dictionary.PollCountMetricName:     &strategy.PollCountStrategy[Num]{},
		dictionary.RandomValueMetricName:   &strategy.RandomValueStrategy[Num]{},
	}

	if _, ok := strategyMap[metricName]; !ok {
		panic("unknown metric type")
	}

	return strategyMap[metricName]
}
