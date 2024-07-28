package manager

import (
	"fmt"
	"github.com/go-resty/resty"
	parserAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/value_handler/abstract_factory"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type MetricManager[Num constraint.Numeric] struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
}

func (metricManager *MetricManager[Num]) CollectAndSend() {
	// todo: Првоерка на nil metricManager.Int64MetricDataCollector.
	var strategy interfaces.ParsingStrategy[Num]
	var metricDTO dto.Metric[Num]
	var collectedMetricData dto.CollectedMetricData

	metricManager.MetricDataCollector.CollectMetricData(&collectedMetricData)

	for _, metricName := range metricManager.MetricList {
		strategy = parserAbstractFactory.GetStrategy[Num](metricName)
		metricDTO = strategy.GetMetric(metricName, collectedMetricData)
		metricManager.sendMetrics(metricDTO)
	}
}

func (metricManager *MetricManager[Num]) UpdateValue(metricDTO dto.Metric[Num]) {
	repos := repository.GetInstance[Num](metricDTO.Type)

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metricDTO, repos)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO)

	repos.UpdateMetric(metricDTO)
}

func (metricManager *MetricManager[Num]) Get(metricDTO dto.Metric[Num]) (dto.Metric[Num], bool) {
	repos := repository.GetInstance[Num](metricDTO.Type)

	metricDTOFromDb, isSet := repos.GetMetric(metricDTO)
	if isSet {
		return metricDTOFromDb, true
	}
	return metricDTO, false
}

//func (metricManager *MetricManager) GetList(metricDTO dto.Metric) []dto.Metric {
//
//}

func (metricManager *MetricManager[Num]) sendMetrics(metricDTO dto.Metric[Num]) {
	request := resty.NewRequest()
	_, err := request.Post(getPreparedURL(metricDTO))
	if err != nil {
		panic(err)
	}
}

func getPreparedURL[Num constraint.Numeric](metricDTO dto.Metric[Num]) string {
	endpoint := "localhost:8080"
	if dictionary.GaugeMetricType == metricDTO.Type {
		return fmt.Sprintf("http://%s/update/%s/%s/%.2f", endpoint, metricDTO.Type, metricDTO.Name, metricDTO.Value)
	} else {
		return fmt.Sprintf("http://%s/update/%s/%s/%d", endpoint, metricDTO.Type, metricDTO.Name, metricDTO.Value)
	}
}
