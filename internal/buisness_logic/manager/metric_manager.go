package manager

import (
	"context"
	"fmt"
	"github.com/go-resty/resty"
	parserAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/value_handler/abstract_factory"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"time"
)

type MetricManager struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
	Repository          repository.Repository
}

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

var (
	endpoint   = "http://localhost:8080"
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

			var strategies = make(map[string]interfaces.ParsingStrategy)
			metricList = []dto.Metric{}
			for _, metricName := range metricManager.MetricList {
				if nil == strategies[metricName] {
					strategies[metricName] = parserAbstractFactory.GetStrategy(metricName)
				}

				strategy := parserAbstractFactory.GetStrategy(metricName)
				metrics := strategy.GetMetric(metricName, metricCollectedData)
				metricList = append(metricList, metrics)
			}

			time.Sleep(pollInterval)
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
				metricManager.sendMetric(metrics)
			}

			time.Sleep(reportInterval)
		}
	}
}

func (metricManager *MetricManager) UpdateValue(metric dto.Metric) {
	repos := repository.GetInstance()

	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metric, repos)
	metric = valueHandler.GetMutatedValueMetric(metric)

	repos.UpdateMetric(metric)
}

func (metricManager *MetricManager) Get(metric dto.Metric) (dto.Metric, bool) {
	repos := repository.GetInstance()

	metricDTOFromDb, isSet := repos.GetMetric(metric)
	if isSet {
		return metricDTOFromDb, true
	}
	return metric, false
}

func (metricManager *MetricManager) sendMetric(metric dto.Metric) {
	_, _ = resty.NewRequest().Post(getPreparedURL(metric))
}

func getPreparedURL(metric dto.Metric) string {
	if metric.Type == dictionary.GaugeMetricType {
		return fmt.Sprintf("%s/update/%s/%s/%.2f", endpoint, metric.Type, metric.Name, metric.Value)
	} else {
		return fmt.Sprintf("%s/update/%s/%s/%d", endpoint, metric.Type, metric.Name, metric.Delta)
	}
}
