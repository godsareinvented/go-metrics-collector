package parser

import (
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"io"
	"net/http"
)

type JsonParser struct{}

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
