package manager

import (
	parserAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/value_handler/abstract_factory"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type MetricManager struct {
	MetricList    []string
	DataCollector interfaces.MetricDataCollectorInterface
	strategies    map[string]interfaces.ParsingStrategyInterface
}

func (metricManager *MetricManager) Collect() []dto.Metrics {
	if nil == metricManager.DataCollector {
		panic("nil DataCollector")
	}

	var metric dto.Metrics
	var metricList []dto.Metrics
	var collectedMetricData dto.CollectedMetricData

	metricManager.DataCollector.CollectMetricData(&collectedMetricData)

	for _, metricName := range metricManager.MetricList {
		metric = metricManager.strategies[metricName].GetMetric(metricName, collectedMetricData)
		metricList = append(metricList, metric)
	}

	return metricList
}

func (metricManager *MetricManager) UpdateMetric(metricDTO dto.Metrics) {
	repos := config.Configuration.Repository

	metricFromStorage, isSet, _ := repos.GetMetric(metricDTO)

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metricDTO)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO, metricFromStorage, isSet)

	_, err := repos.UpdateMetric(metricDTO)
	if nil != err {
		// todo: Надо пересмотреть выплёвывание ошибок.
		panic("Error updating metric: " + err.Error())
	}

	if 0 == config.Configuration.StoreInterval {
		_ = metricManager.ExportTo(config.Configuration.PermanentStorage)
	}
}

func (metricManager *MetricManager) ImportFrom(permanentStorage *interfaces.PermanentStorage) error {
	metricList, err := (*permanentStorage).Import()
	if nil != err {
		return err
	}

	for _, metric := range metricList {
		metricManager.UpdateMetric(metric)
	}

	return nil
}

func (metricManager *MetricManager) ExportTo(permanentStorage *interfaces.PermanentStorage) error {
	metricList, err := config.Configuration.Repository.GetAllMetrics()
	if nil != err {
		return err
	}

	err = (*permanentStorage).Export(metricList)
	if err != nil {
		return err
	}

	return nil
}

func (metricManager *MetricManager) Init() {
	metricManager.initStrategyList()
}

func (metricManager *MetricManager) initStrategyList() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategyInterface)

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}
