package manager

import (
	"context"
	parserAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/value_handler/abstract_factory"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/sender"
	"time"
)

type MetricManager struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
	strategies          map[string]interfaces.ParsingStrategy
	sender              *sender.MetricSender
}

var (
	metricList []dto.Metric
)

func (metricManager *MetricManager) CollectAndSend(ctx context.Context) {
	if metricManager.MetricList == nil {
		panic("metric list is empty")
	}

	go metricManager.collect(ctx)
	go metricManager.send(ctx)

	select {
	case <-ctx.Done():
		return
	}
}

func (metricManager *MetricManager) collect(ctx context.Context) {
	var metricCollectedData dto.CollectedMetricData
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricManager.MetricDataCollector.CollectMetricData(&metricCollectedData)

			metricList = []dto.Metric{}
			for _, metricName := range metricManager.MetricList {
				strategy := metricManager.strategies[metricName]
				metrics := strategy.GetMetric(metricName, metricCollectedData)
				metricList = append(metricList, metrics)
			}

			time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
		}
	}
}

func (metricManager *MetricManager) send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, metrics := range metricList {
				metricManager.sender.Send(metrics)
			}

			time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
		}
	}
}

func (metricManager *MetricManager) UpdateValue(metric dto.Metric) {
	repos := config.Configuration.Repository

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metric, repos)
	metric = valueHandler.GetMutatedValueMetric(metric)

	repos.UpdateMetric(metric)
}

func (metricManager *MetricManager) Get(metric dto.Metric) (dto.Metric, bool) {
	repos := config.Configuration.Repository

	metricDTOFromDb, isSet := repos.GetMetric(metric)
	if isSet {
		return metricDTOFromDb, true
	}
	return metric, false
}

func (metricManager *MetricManager) GetList() []dto.Metric {
	repos := config.Configuration.Repository

	return repos.GetAllMetrics()
}

func (metricManager *MetricManager) Init() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategy)
	metricManager.sender = sender.NewSender()

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
	}
}
