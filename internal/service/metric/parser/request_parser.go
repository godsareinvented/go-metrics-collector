package parser

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"net/http"
)

type RequestParser struct{}

func (parser *RequestParser) GetParsedMetric(request *http.Request) dto.Metric {
	return dto.Metric{
		Type:  request.PathValue("type"),
		Name:  request.PathValue("name"),
		Value: request.PathValue("value"),
	}
}
