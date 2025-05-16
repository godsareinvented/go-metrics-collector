package manager

import (
	"context"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/data_collector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	errorRate = 0.7
)

var (
	processedMetricList []dto.Metric
	mu                  sync.Mutex
	requestCount        = atomic.Uint32{}
)

// TestCollectAndSend Тест будет работать при условии, что значение reportInterval - минимум, 1 секунда,
// т.к. за это время все метрики должны успеть уйти на сервер
func TestCollectAndSend(t *testing.T) {
	server := httptest.NewServer(router(t))
	defer server.Close()

	oldEndpoint := endpoint
	endpoint = server.URL
	defer func() {
		endpoint = oldEndpoint
	}()

	metricManager := MetricManager{MetricDataCollector: &data_collector.MetricDataCollector{}, MetricList: dictionary.MetricNameList[:]}

	ctx, cancel := context.WithTimeout(context.Background(), reportInterval+1*time.Second)
	defer cancel()
	metricManager.CollectAndSend(ctx)

	testMetricList(t, "test collect method", metricList)
	testMetricList(t, "test send method", processedMetricList)
}

func testMetricList(t *testing.T, testName string, metrics []dto.Metric) {
	t.Run(testName, func(t *testing.T) {
		metricCountMap := make(map[string]int)
		zeroValueMetricCount := 0
		allowedMetricTypes := []string{dictionary.GaugeMetricType, dictionary.CounterMetricType}

		for _, metric := range metrics {
			if _, ok := metricCountMap[metric.Name]; !ok {
				metricCountMap[metric.Name] = 0
			}

			metricCountMap[metric.Name]++

			if metricCountMap[metric.Name] > 1 {
				t.Fatalf("metric '%s' passed more than once", metric.Name)
			}

			require.Containsf(t, allowedMetricTypes, metric.Type, "metric %s is of a type not allowed", metric.Name)
			require.Containsf(t, dictionary.MetricNameList, metric.Name, "metric %s is of a name not allowed", metric.Name)

			if metric.Type == dictionary.GaugeMetricType {
				require.GreaterOrEqualf(t, metric.Value, 0.0, "%s metric value must be non-negative", metric.Name)
				if metric.Value == 0.0 {
					zeroValueMetricCount++
				}
				continue
			}

			require.GreaterOrEqualf(t, metric.Delta, int64(0), "%s metric value must be non-negative", metric.Name)
			if metric.Value == 0 {
				zeroValueMetricCount++
			}
		}

		require.Equalf(t, len(dictionary.MetricNameList), len(processedMetricList), "Passed %d metrics, need %d", len(processedMetricList), len(dictionary.MetricNameList))

		zeroValueMetricsPercent := float64(zeroValueMetricCount) / float64(len(dictionary.MetricNameList))
		require.LessOrEqualf(t, zeroValueMetricsPercent, errorRate, "Too many metrics with zero values (%.0f%%)", zeroValueMetricsPercent*100)
	})
}

func router(t *testing.T) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/update/{type}/{name}/{value}", func(_ http.ResponseWriter, r *http.Request) {
		testName := nameOfTestByMetricName(r.PathValue("name"))
		t.Run(testName, handle(r))
	})

	return r
}

func handle(r *http.Request) func(t *testing.T) {
	return func(t *testing.T) {
		if !assert.Equal(t, http.MethodPost, r.Method) {
			return
		}

		MType, MName, MValue := parsedMetricValues(r)
		if !assert.NotNil(t, MType) || !assert.NotNil(t, MName) || !assert.NotNil(t, MValue) {
			return
		}

		metrics := dto.Metric{Type: MType, Name: MName}
		var err error
		if MType == dictionary.GaugeMetricType {
			metrics.Value, err = strconv.ParseFloat(MValue, 64)
		} else {
			metrics.Delta, err = strconv.ParseInt(MValue, 10, 64)
		}

		require.NoErrorf(t, err, "%s metric value should be of correct type", MName)

		mu.Lock()
		processedMetricList = append(processedMetricList, metrics)
		mu.Unlock()
	}
}

func nameOfTestByMetricName(MName string) string {
	requestCount.Swap(requestCount.Load() + 1)

	if MName != "" {
		return fmt.Sprintf("request #%d (%s metric)", requestCount.Load(), MName)
	}
	return fmt.Sprintf("request #%d", requestCount.Load())
}

func parsedMetricValues(r *http.Request) (string, string, string) {
	return r.PathValue("type"),
		r.PathValue("name"),
		r.PathValue("value")
}
