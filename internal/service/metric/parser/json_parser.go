package parser

import (
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"io"
	"net/http"
)

type JsonParser struct{}

// GetMetricDTO todo: Убрать в будущем значения по умолчанию для целочисленного и вещественного значений.
func (jp *JsonParser) GetMetricDTO(request *http.Request) (dto.Metrics, error) {
	var intVal int64 = 0
	var floatVal = 0.0
	var metricDTO = dto.Metrics{
		Delta: &intVal,
		Value: &floatVal,
	}

	body, err := io.ReadAll(request.Body)

	if err != nil {
		return metricDTO, err
	}

	err = json.Unmarshal(body, &metricDTO)

	if err != nil {
		return metricDTO, err
	}

	return metricDTO, nil
}

func (jp *JsonParser) GetMetricBatch(request *http.Request) ([]dto.Metrics, error) {
	var metrics []dto.Metrics

	body, err := io.ReadAll(request.Body)

	if err != nil {
		return metrics, err
	}

	err = json.Unmarshal(body, &metrics)

	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
