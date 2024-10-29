package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type RandomValueStrategy struct{}

func (strategy *RandomValueStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = metricData.RandomValue
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
