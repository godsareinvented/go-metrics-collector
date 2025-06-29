package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
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
				strategy:   &AllocStrategy{},
				want:       generaldto.Metrics{ID: dictionary.AllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.BuckHashSysMetricName,
				strategy:   &BuckHashSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.BuckHashSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.FreesMetricName,
				strategy:   &FreesStrategy{},
				want:       generaldto.Metrics{ID: dictionary.FreesMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.GCCPUFractionMetricName,
				strategy:   &GCCPUFractionStrategy{},
				want:       generaldto.Metrics{ID: dictionary.GCCPUFractionMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.GCSysMetricName,
				strategy:   &GCSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.GCSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapAllocMetricName,
				strategy:   &HeapAllocStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapAllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapIdleMetricName,
				strategy:   &HeapIdleStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapIdleMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapInuseMetricName,
				strategy:   &HeapInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapObjectsMetricName,
				strategy:   &HeapObjectsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapObjectsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapReleasedMetricName,
				strategy:   &HeapReleasedStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapReleasedMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapSysMetricName,
				strategy:   &HeapSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.HeapSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.LastGCMetricName,
				strategy:   &LastGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.LastGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.LookupsMetricName,
				strategy:   &LookupsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.LookupsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MCacheInuseMetricName,
				strategy:   &MCacheInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MCacheInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MCacheSysMetricName,
				strategy:   &MCacheSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MCacheSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MSpanInuseMetricName,
				strategy:   &MSpanInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MSpanInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MSpanSysMetricName,
				strategy:   &MSpanSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MSpanSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MallocsMetricName,
				strategy:   &MallocsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.MallocsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NextGCMetricName,
				strategy:   &NextGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.NextGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NumForcedGCMetricName,
				strategy:   &NumForcedGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.NumForcedGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NumGCMetricName,
				strategy:   &NumGCStrategy{},
				want:       generaldto.Metrics{ID: dictionary.NumGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.OtherSysMetricName,
				strategy:   &OtherSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.OtherSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.PauseTotalNsMetricName,
				strategy:   &PauseTotalNsStrategy{},
				want:       generaldto.Metrics{ID: dictionary.PauseTotalNsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.StackInuseMetricName,
				strategy:   &StackInuseStrategy{},
				want:       generaldto.Metrics{ID: dictionary.StackInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.StackSysMetricName,
				strategy:   &StackSysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.StackSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.SysMetricName,
				strategy:   &SysStrategy{},
				want:       generaldto.Metrics{ID: dictionary.SysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.TotalAllocMetricName,
				strategy:   &TotalAllocStrategy{},
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
			test.strategy.GetMetric(&metric, &collectedData)
			assert.Equal(t, test.want, metric)
		}
	})
}
