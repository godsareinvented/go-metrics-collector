package action

import (
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/value_handler/factory"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

func UpdateValue(metricDTO dto.Metric) {
	valueHandler := factory.GetValueHandler(metricDTO)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO)

	repository.MetricRepository.UpdateMetric(metricDTO)
}
