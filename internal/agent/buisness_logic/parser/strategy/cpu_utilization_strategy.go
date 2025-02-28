package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type CPUutilizationStrategy struct {
	logicalCpuNumber uint
	metricName       string
}

func (strategy *CPUutilizationStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	if !isLogicalCpuPercentsValid(&metricData.CPUPercentList) {
		return ErrEmptyCollectedData
	}

	var value = metricData.CPUPercentList[strategy.logicalCpuNumber]
	metric.MType = dictionary.GaugeMetricType
	metric.MName = strategy.metricName
	metric.Value = &value

	return nil
}

func NewCPUutilizationStrategy(logicalCpuNumber uint, metricName string) interfaces.ParsingStrategyInterface {
	return &CPUutilizationStrategy{
		logicalCpuNumber: logicalCpuNumber,
		metricName:       metricName,
	}
}
