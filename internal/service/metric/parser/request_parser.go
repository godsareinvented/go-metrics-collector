package parser

import (
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"net/http"
	"strconv"
)

type RequestParser struct{}

func (rp *RequestParser) GetMetricDTO(request *http.Request, parsingValueFlag bool) (dto.Metrics, error) {
	metricType, metricName, metricValue := getParsedRequest(request)

	var intVal int64
	var floatVal float64
	var err error

	if parsingValueFlag {
		intVal, err = getParsedDelta(metricType, metricValue)
		if nil != err {
			return dto.Metrics{}, err
		}
		floatVal, err = getParsedValue(metricType, metricValue)
		if nil != err {
			return dto.Metrics{}, err
		}
	}

	return dto.Metrics{
		MType: metricType,
		MName: metricName,
		Delta: &intVal,
		Value: &floatVal,
	}, nil
}

func getParsedRequest(request *http.Request) (string, string, string) {
	return chi.URLParam(request, "type"),
		chi.URLParam(request, "name"),
		chi.URLParam(request, "value")
}

func getParsedDelta(metricType string, metricValue string) (int64, error) {
	if dictionary.CounterMetricType != metricType {
		return 0, nil
	}

	return strconv.ParseInt(metricValue, 10, 64)
}

func getParsedValue(metricType string, metricValue string) (float64, error) {
	if dictionary.GaugeMetricType != metricType {
		return 0.0, nil
	}

	return strconv.ParseFloat(metricValue, 64)
}
