package manager

import (
	"context"
	"fmt"
	parserAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/parser/abstract_factory"
	valueHandlerAbstractFactory "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/value_handler/abstract_factory"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
	"net/http"
	"reflect"
	"time"
)

type MetricManager struct {
	MetricDataCollector data_collector.MetricDataCollector
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
	go metricManager.collect(ctx)
	go metricManager.send(ctx)

	select {
	case <-ctx.Done():
		return
	}
}

func (metricManager *MetricManager) collect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricCollectedData := metricManager.MetricDataCollector.GetMetricData()

			metricList = []dto.Metric{}
			for _, metricName := range dictionary.MetricNameList {
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
				metricManager.sendMetrics(metrics)
			}

			time.Sleep(reportInterval)
		}
	}
}

func (metricManager *MetricManager) UpdateValue(metricDTO dto.Metric) {
	valueHandler := valueHandlerAbstractFactory.GetValueHandler(metricDTO)
	metricDTO = valueHandler.GetMutatedValueMetric(metricDTO)

	repository.MetricRepository.UpdateMetric(metricDTO)
}

func (metricManager *MetricManager) sendMetrics(metricDTO dto.Metric) {
	response, err := http.Post(getPreparedURL(metricDTO), "text/plain", http.NoBody)
	if err != nil {
		return
	}

	_ = response.Body.Close()
}

func getPreparedURL(metricDTO dto.Metric) string {
	if reflect.TypeOf(metricDTO.Value).Name() == "float64" {
		preparedMetricValue := metricDTO.Value.(float64)
		return fmt.Sprintf("%s/update/%s/%s/%.2f", endpoint, metricDTO.Type, metricDTO.Name, preparedMetricValue)
	} else {
		preparedMetricValue := metricDTO.Value.(int64)
		return fmt.Sprintf("%s/update/%s/%s/%d", endpoint, metricDTO.Type, metricDTO.Name, preparedMetricValue)
	}
}
