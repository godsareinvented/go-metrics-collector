package manager

import (
	"fmt"
	"github.com/go-resty/resty"
	parserAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/value_handler/abstract_factory"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type MetricManager struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
}

func (metricManager *MetricManager) CollectAndSend() {
	if nil == metricManager.MetricList {
		panic("metric list is empty")
	}

	var strategies = make(map[string]interfaces.ParsingStrategy)
	var metricDTO dto.Metric
	var collectedMetricData dto.CollectedMetricData

	metricManager.MetricDataCollector.CollectMetricData(&collectedMetricData)

	for _, metricName := range metricManager.MetricList {
		if nil == strategies[metricName] {
			strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
		}
		metricDTO = strategies[metricName].GetMetric(metricName, collectedMetricData)
		metricManager.sendMetrics(metricDTO)
	}
}

func (metricManager *MetricManager) UpdateValue(metricDTO dto.Metric) {
	repos := repository.GetInstance()

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metricDTO, repos)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO)

	repos.UpdateMetric(metricDTO)
}

func (metricManager *MetricManager) Get(metricDTO dto.Metric) (dto.Metric, bool) {
	repos := repository.GetInstance()

	metricDTOFromDb, isSet := repos.GetMetric(metricDTO)
	if isSet {
		return metricDTOFromDb, true
	}
	return metricDTO, false
}

func (metricManager *MetricManager) GetList() []dto.Metric {
	repos := repository.GetInstance()

	return repos.GetAllMetrics()
}

func (metricManager *MetricManager) sendMetrics(metricDTO dto.Metric) {
	request := resty.NewRequest()
	_, err := request.Post(getPreparedURL(metricDTO))
	if err != nil {
		panic(err)
	}
}

func getPreparedURL(metricDTO dto.Metric) string {
	endpoint := "localhost:8080"
	if dictionary.GaugeMetricType == metricDTO.Type {
		return fmt.Sprintf("http://%s/update/%s/%s/%.2f", endpoint, metricDTO.Type, metricDTO.Name, metricDTO.Value)
	} else {
		return fmt.Sprintf("http://%s/update/%s/%s/%d", endpoint, metricDTO.Type, metricDTO.Name, metricDTO.Delta)
	}
}
