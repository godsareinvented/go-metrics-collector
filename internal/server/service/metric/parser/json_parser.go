package parser

import (
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"net/http"
)

type JsonParser struct{}

// GetMetricDTO todo: Убрать в будущем значения по умолчанию для целочисленного и вещественного значений.
func (jp *JsonParser) GetMetricDTO(request *http.Request) (dto.Metrics, error) {
	var metric = dto.Metrics{}

	err := json.NewDecoder(request.Body).Decode(&metric)
	if nil != err {
		return metric, err
	}

	return metric, nil
}

func (jp *JsonParser) GetMetricBatch(request *http.Request) ([]dto.Metrics, error) {
	var metrics []dto.Metrics

	err := json.NewDecoder(request.Body).Decode(&metrics)
	if nil != err {
		return nil, err
	}

	return metrics, nil
}
