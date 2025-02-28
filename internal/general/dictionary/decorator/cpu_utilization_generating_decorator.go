package decorator

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"strconv"
)

func AddCpuUtilizationMetricNames(metricList []string, logicalCpuCount int) ([]string, error) {
	var metricName string
	for i := 0; i < logicalCpuCount; i++ {
		metricName = dictionary.CPUutilizationMetricName + strconv.Itoa(i)
		metricList = append(metricList, metricName)
	}

	return metricList, nil
}
