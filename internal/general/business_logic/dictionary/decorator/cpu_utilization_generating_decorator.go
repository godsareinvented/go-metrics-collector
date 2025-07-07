package decorator

import (
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"strconv"
)

func AddCpuUtilizationMetricNames(metricList []string, logicalCpuCount int) ([]string, error) {
	var metricName string
	for i := 0; i < logicalCpuCount; i++ {
		metricName = dictionary.CpuUtilizationMetricName + strconv.Itoa(i)
		metricList = append(metricList, metricName)
	}

	return metricList, nil
}
