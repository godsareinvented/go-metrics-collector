package cpu

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type CpuUtilizationStrategy struct {
	logicalCpuNumber uint
	metricName       string
}

func (strategy *CpuUtilizationStrategy) FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData) {
	metric.ID = strategy.metricName
	metric.MType = dictionary.GaugeMetricType
	metric.Delta = nil
	metric.Value = &metricData.CPUPercentList[strategy.logicalCpuNumber]
}

func NewCpuUtilizationStrategy(logicalCpuNumber uint, metricName string) interfaces.ParsingStrategy {
	return &CpuUtilizationStrategy{
		logicalCpuNumber: logicalCpuNumber,
		metricName:       metricName,
	}
}
