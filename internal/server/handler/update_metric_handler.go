package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
	"net/http"
	"strings"
)

func UpdateMetric(responseWriter http.ResponseWriter, request *http.Request) {
	// todo: Удалить и перепрвоерить перед коммитом.
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if chi.URLParam(request, "type") == dictionary.GaugeMetricType {
		handle[float64](responseWriter, request)
		return
	}
	handle[int64](responseWriter, request)
}

func handle[Num constraint.Numeric](responseWriter http.ResponseWriter, request *http.Request) {
	requestParser := parser.RequestParser[Num]{}
	metricDTO := requestParser.GetMetricDTO(request)

	err := validator.New().Struct(metricDTO)

	if nil != err {
		message, statusCode := processError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	metricRepository := repository.NewInstance[Num](mem_storage.NewInstance())
	metricManager := manager.MetricManager[Num]{Repository: metricRepository}
	metricManager.UpdateValue(metricDTO)
}

func processError(error error) (string, int) {
	errors := error.(validator.ValidationErrors)
	if strings.Contains(errors[0].Field(), "Name") {
		return "Metric name not passed on or incorrect", http.StatusNotFound
	}
	return "incorrect metric data", http.StatusBadRequest
}
