package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
)

type InputMetrics struct {
	ID    string   `json:"id"               validate:"required,omitempty,required"`
	MType string   `json:"type"             validate:"required,contains=gauge|contains=counter"`
	MName string   `json:"name"             validate:""`
	Delta *int64   `json:"delta,omitempty"  validate:""`
	Value *float64 `json:"value,omitempty"  validate:""`
}

func GetMetricJson(responseWriter http.ResponseWriter, request *http.Request) {
	requestParser := parser.JsonParser{}
	metric, err := requestParser.GetMetricDTO(request)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	InputMetric := InputMetrics{
		ID:    metric.ID,
		MType: metric.MType,
	}
	err = validator.New().Struct(InputMetric)
	if nil != err {
		message, statusCode := ProcessValidationError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	metricManager := manager.MetricManager{}
	resultingMetric, isSet := metricManager.GetByID(metric)

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
	responseWriter.WriteHeader(http.StatusOK)
	return
}
