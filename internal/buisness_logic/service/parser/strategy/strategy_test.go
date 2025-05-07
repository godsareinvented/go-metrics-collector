package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	tests = []struct {
		metricName string
		strategy   interfaces.ParsingStrategy
		want       dto.Metric
	}{
		{
			metricName: dictionary.AllocMetricName,
			strategy:   &AllocStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.AllocMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.BuckHashSysMetricName,
			strategy:   &BuckHashSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.BuckHashSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.FreesMetricName,
			strategy:   &FreesStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.FreesMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.GCCPUFractionMetricName,
			strategy:   &GCCPUFractionStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.GCCPUFractionMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.GCSysMetricName,
			strategy:   &GCSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.GCSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.HeapAllocMetricName,
			strategy:   &HeapAllocStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.HeapAllocMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.HeapIdleMetricName,
			strategy:   &HeapIdleStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.HeapIdleMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.HeapInuseMetricName,
			strategy:   &HeapInuseStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.HeapInuseMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.HeapObjectsMetricName,
			strategy:   &HeapObjectsStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.HeapObjectsMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.HeapReleasedMetricName,
			strategy:   &HeapReleasedStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.HeapReleasedMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.HeapSysMetricName,
			strategy:   &HeapSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.HeapSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.LastGCMetricName,
			strategy:   &LastGCStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.LastGCMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.LookupsMetricName,
			strategy:   &LookupsStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.LookupsMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.MCacheInuseMetricName,
			strategy:   &MCacheInuseStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.MCacheInuseMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.MCacheSysMetricName,
			strategy:   &MCacheSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.MCacheSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.MSpanInuseMetricName,
			strategy:   &MSpanInuseStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.MSpanInuseMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.MSpanSysMetricName,
			strategy:   &MSpanSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.MSpanSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.MallocsMetricName,
			strategy:   &MallocsStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.MallocsMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.NextGCMetricName,
			strategy:   &NextGCStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.NextGCMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.NumForcedGCMetricName,
			strategy:   &NumForcedGCStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.NumForcedGCMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.NumGCMetricName,
			strategy:   &NumGCStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.NumGCMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.OtherSysMetricName,
			strategy:   &OtherSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.OtherSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.PauseTotalNsMetricName,
			strategy:   &PauseTotalNsStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.PauseTotalNsMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.StackInuseMetricName,
			strategy:   &StackInuseStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.StackInuseMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.StackSysMetricName,
			strategy:   &StackSysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.StackSysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.SysMetricName,
			strategy:   &SysStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.SysMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.TotalAllocMetricName,
			strategy:   &TotalAllocStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.TotalAllocMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.RandomValueMetricName,
			strategy:   &RandomValueStrategy{},
			want:       dto.Metric{Type: dictionary.GaugeMetricType, Name: dictionary.RandomValueMetricName, Value: 0.0},
		},
		{
			metricName: dictionary.PollCountMetricName,
			strategy:   &PollCountStrategy{},
			want:       dto.Metric{Type: dictionary.CounterMetricType, Name: dictionary.PollCountMetricName, Value: int64(0)},
		},
	}
)

func TestGetMetric(t *testing.T) {
	t.Run("strategy get metric", func(t *testing.T) {
		for _, test := range tests {
			metrics := test.strategy.GetMetric(test.metricName, dto.CollectedMetricData{})
			assert.Equal(t, test.want, metrics)
		}
	})
}
