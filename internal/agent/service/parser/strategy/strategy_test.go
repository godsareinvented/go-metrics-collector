package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/service/parser/strategy/mem"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMetric(t *testing.T) {
	var (
		delta int64 = 0
		value       = 0.0
		tests       = []struct {
			metricName string
			strategy   interfaces.ParsingStrategy
			want       generaldto.Metrics
		}{
			{
				metricName: dictionary.AllocMetricName,
				strategy:   &mem.AllocStrategy{},
				want:       generaldto.Metrics{ID: dictionary.AllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.BuckHashSysMetricName,
				strategy:   &mem.BuckHashSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.BuckHashSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.FreesMetricName,
				strategy:   &mem.FreesStrategy{},
				want:       generaldto.Metrics{ID: dictionary.FreesMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.GCCPUFractionMetricName,
				strategy:   &mem.GCCPUFractionStrategy{},
				want:       generaldto.Metrics{ID: dictionary.GCCPUFractionMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.GCSysMetricName,
				strategy:   &mem.GCSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.GCSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapAllocMetricName,
				strategy:   &mem.HeapAllocStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapAllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapIdleMetricName,
				strategy:   &mem.HeapIdleStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapIdleMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapInuseMetricName,
				strategy:   &mem.HeapInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapObjectsMetricName,
				strategy:   &mem.HeapObjectsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapObjectsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapReleasedMetricName,
				strategy:   &mem.HeapReleasedStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapReleasedMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapSysMetricName,
				strategy:   &mem.HeapSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.LastGCMetricName,
				strategy:   &mem.LastGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.LastGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.LookupsMetricName,
				strategy:   &mem.LookupsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.LookupsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MCacheInuseMetricName,
				strategy:   &mem.MCacheInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MCacheInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MCacheSysMetricName,
				strategy:   &mem.MCacheSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MCacheSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MSpanInuseMetricName,
				strategy:   &mem.MSpanInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MSpanInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MSpanSysMetricName,
				strategy:   &mem.MSpanSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MSpanSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MallocsMetricName,
				strategy:   &mem.MallocsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MallocsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NextGCMetricName,
				strategy:   &mem.NextGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.NextGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NumForcedGCMetricName,
				strategy:   &mem.NumForcedGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.NumForcedGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NumGCMetricName,
				strategy:   &mem.NumGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.NumGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.OtherSysMetricName,
				strategy:   &mem.OtherSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.OtherSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.PauseTotalNsMetricName,
				strategy:   &mem.PauseTotalNsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.PauseTotalNsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.StackInuseMetricName,
				strategy:   &mem.StackInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.StackInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.StackSysMetricName,
				strategy:   &mem.StackSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.StackSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.SysMetricName,
				strategy:   &mem.SysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.SysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.TotalAllocMetricName,
				strategy:   &mem.TotalAllocStrategy{},
				want:       generaldto.Metrics{ID: dictionary.TotalAllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.RandomValueMetricName,
				strategy:   &RandomValueStrategy{},
				want:       generaldto.Metrics{ID: dictionary.RandomValueMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.PollCountMetricName,
				strategy:   &PollCountStrategy{},
				want:       generaldto.Metrics{ID: dictionary.PollCountMetricName, MType: dictionary.CounterMetricType, Value: nil, Delta: &delta},
			},
		}
	)

	t.Run("strategy get metric", func(t *testing.T) {
		var metric generaldto.Metrics
		var collectedData agentdto.CollectedMetricData

		for _, test := range tests {
			test.strategy.FillMetric(&metric, &collectedData)
			assert.Equal(t, test.want, metric)
		}
	})
}
