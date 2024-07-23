package manager

import (
	"context"
	"fmt"
	"github.com/go-resty/resty"
	parserFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/parser/factory"
	valueHandlerFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/value_handler/factory"
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"reflect"
	"time"
)

type MetricManager[Num constraint.Numeric] struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
	Repository          repository.Repository[Num]
}

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

var (
	endpoint          = "http://localhost:8080"
	int64MetricList   []dto.Metric[int64]
	float64MetricList []dto.Metric[float64]
)

func (metricManager *MetricManager[Num]) CollectAndSend(ctx context.Context) {
	go metricManager.collect(ctx)
	go metricManager.send(ctx)

	select {
	case <-ctx.Done():
		return
	}
}

func (metricManager *MetricManager[Num]) collect(ctx context.Context) {
	var metricCollectedData dto.CollectedMetricData
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricManager.MetricDataCollector.CollectMetricData(&metricCollectedData)

			if reflect.TypeFor[Num]().Name() == "int64" {
				collectCounterMetrics(metricCollectedData)
			} else {
				collectGaugeMetrics(metricCollectedData)
			}

			time.Sleep(pollInterval)
		}
	}
}

func (metricManager *MetricManager[Num]) send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, metrics := range int64MetricList {
				_, _ = resty.NewRequest().Post(getPreparedURL[int64](metrics))
			}

			for _, metrics := range float64MetricList {
				_, _ = resty.NewRequest().Post(getPreparedURL[float64](metrics))
			}

			time.Sleep(reportInterval)
		}
	}
}

func (metricManager *MetricManager[Num]) UpdateValue(metric dto.Metric[Num]) {
	repos := repository.GetInstance(metric)

	valueHandler := valueHandlerFactory.GetValueHandler(metric, repos)
	metric = valueHandler.GetMutatedValueMetric(metric)

	repos.UpdateMetric(metric)
}

func (metricManager *MetricManager[Num]) Get(metric dto.Metric[Num]) (dto.Metric[Num], bool) {
	repos := repository.GetInstance(metric)

	metricDTOFromDb, isSet := repos.GetMetric(metric)
	if isSet {
		return metricDTOFromDb, true
	}
	return metric, false
}

// todo: Временное решение.
func (metricManager *MetricManager[Num]) _(metric dto.Metric[Num]) {
	_, _ = resty.NewRequest().Post(getPreparedURL(metric))
}

func getPreparedURL[Num constraint.Numeric](metric dto.Metric[Num]) string {
	if metric.Type == dictionary.GaugeMetricType {
		return fmt.Sprintf("%s/update/%s/%s/%.2f", endpoint, metric.Type, metric.Name, metric.Value)
	} else {
		return fmt.Sprintf("%s/update/%s/%s/%d", endpoint, metric.Type, metric.Name, metric.Value)
	}
}

func collectCounterMetrics(metricCollectedData dto.CollectedMetricData) {
	int64MetricList = []dto.Metric[int64]{}
	for _, metricName := range dictionary.CounterMetricNameList {
		strategy := parserFactory.GetStrategy[int64](metricName)
		metrics := strategy.GetMetric(metricName, metricCollectedData)
		int64MetricList = append(int64MetricList, metrics)
	}
}

func collectGaugeMetrics(metricCollectedData dto.CollectedMetricData) {
	float64MetricList = []dto.Metric[float64]{}
	for _, metricName := range dictionary.GaugeMetricNameList {
		strategy := parserFactory.GetStrategy[float64](metricName)
		metrics := strategy.GetMetric(metricName, metricCollectedData)
		float64MetricList = append(float64MetricList, metrics)
	}
}
