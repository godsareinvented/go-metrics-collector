package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
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
			want       dto.Metrics
		}{
			{
				metricName: dictionary.AllocMetricName,
				strategy:   &AllocStrategy{},
				want:       dto.Metrics{ID: dictionary.AllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.BuckHashSysMetricName,
				strategy:   &BuckHashSysStrategy{},
				want:       dto.Metrics{ID: dictionary.BuckHashSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.FreesMetricName,
				strategy:   &FreesStrategy{},
				want:       dto.Metrics{ID: dictionary.FreesMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.GCCPUFractionMetricName,
				strategy:   &GCCPUFractionStrategy{},
				want:       dto.Metrics{ID: dictionary.GCCPUFractionMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.GCSysMetricName,
				strategy:   &GCSysStrategy{},
				want:       dto.Metrics{ID: dictionary.GCSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapAllocMetricName,
				strategy:   &HeapAllocStrategy{},
				want:       dto.Metrics{ID: dictionary.HeapAllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapIdleMetricName,
				strategy:   &HeapIdleStrategy{},
				want:       dto.Metrics{ID: dictionary.HeapIdleMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapInuseMetricName,
				strategy:   &HeapInuseStrategy{},
				want:       dto.Metrics{ID: dictionary.HeapInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapObjectsMetricName,
				strategy:   &HeapObjectsStrategy{},
				want:       dto.Metrics{ID: dictionary.HeapObjectsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapReleasedMetricName,
				strategy:   &HeapReleasedStrategy{},
				want:       dto.Metrics{ID: dictionary.HeapReleasedMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.HeapSysMetricName,
				strategy:   &HeapSysStrategy{},
				want:       dto.Metrics{ID: dictionary.HeapSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.LastGCMetricName,
				strategy:   &LastGCStrategy{},
				want:       dto.Metrics{ID: dictionary.LastGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.LookupsMetricName,
				strategy:   &LookupsStrategy{},
				want:       dto.Metrics{ID: dictionary.LookupsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MCacheInuseMetricName,
				strategy:   &MCacheInuseStrategy{},
				want:       dto.Metrics{ID: dictionary.MCacheInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MCacheSysMetricName,
				strategy:   &MCacheSysStrategy{},
				want:       dto.Metrics{ID: dictionary.MCacheSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MSpanInuseMetricName,
				strategy:   &MSpanInuseStrategy{},
				want:       dto.Metrics{ID: dictionary.MSpanInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MSpanSysMetricName,
				strategy:   &MSpanSysStrategy{},
				want:       dto.Metrics{ID: dictionary.MSpanSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.MallocsMetricName,
				strategy:   &MallocsStrategy{},
				want:       dto.Metrics{ID: dictionary.MallocsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NextGCMetricName,
				strategy:   &NextGCStrategy{},
				want:       dto.Metrics{ID: dictionary.NextGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NumForcedGCMetricName,
				strategy:   &NumForcedGCStrategy{},
				want:       dto.Metrics{ID: dictionary.NumForcedGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.NumGCMetricName,
				strategy:   &NumGCStrategy{},
				want:       dto.Metrics{ID: dictionary.NumGCMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.OtherSysMetricName,
				strategy:   &OtherSysStrategy{},
				want:       dto.Metrics{ID: dictionary.OtherSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.PauseTotalNsMetricName,
				strategy:   &PauseTotalNsStrategy{},
				want:       dto.Metrics{ID: dictionary.PauseTotalNsMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.StackInuseMetricName,
				strategy:   &StackInuseStrategy{},
				want:       dto.Metrics{ID: dictionary.StackInuseMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.StackSysMetricName,
				strategy:   &StackSysStrategy{},
				want:       dto.Metrics{ID: dictionary.StackSysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.SysMetricName,
				strategy:   &SysStrategy{},
				want:       dto.Metrics{ID: dictionary.SysMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.TotalAllocMetricName,
				strategy:   &TotalAllocStrategy{},
				want:       dto.Metrics{ID: dictionary.TotalAllocMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.RandomValueMetricName,
				strategy:   &RandomValueStrategy{},
				want:       dto.Metrics{ID: dictionary.RandomValueMetricName, MType: dictionary.GaugeMetricType, Value: &value, Delta: nil},
			},
			{
				metricName: dictionary.PollCountMetricName,
				strategy:   &PollCountStrategy{},
				want:       dto.Metrics{ID: dictionary.PollCountMetricName, MType: dictionary.CounterMetricType, Value: nil, Delta: &delta},
			},
		}
	)

	t.Run("strategy get metric", func(t *testing.T) {
		for _, test := range tests {
			metrics := test.strategy.GetMetric(test.metricName, dto.CollectedMetricData{})
			assert.Equal(t, test.want, metrics)
		}
	})
}
