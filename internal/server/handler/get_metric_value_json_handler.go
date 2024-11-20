package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
)

type InputMetrics struct {
	ID    string `json:"id"   validate:"required,omitempty,required"`
	MType string `json:"type" validate:"required,contains=gauge|contains=counter"`
}

func GetMetricJson(responseWriter http.ResponseWriter, request *http.Request) {
	requestParser := parser.JsonParser{}
	metric, err := requestParser.GetMetricDTO(request)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	inputMetric := InputMetrics{
		ID:    metric.ID,
		MType: metric.MType,
	}
	err = validator.New().Struct(inputMetric)
	if nil != err {
		message, statusCode := ProcessValidationError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	searchMetric := dto.Metrics{
		ID:    metric.ID,
		MType: metric.MType,
	}
	resultingMetric, isSet, _ := config.Configuration.Repository.GetMetricByID(searchMetric)

	if !isSet {
		http.NotFound(responseWriter, request)
		return
	}

	metricJson, err := json.Marshal(resultingMetric)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	responseWriter.Write(metricJson)
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	return
}
