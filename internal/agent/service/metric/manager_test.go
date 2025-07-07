package metric

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/business_logic/config"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/business_logic/metric/data_collector"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	errorRate = 0.7
)

var (
	processedMetricList []dto.Metrics
	mu                  sync.Mutex
	requestCount        = atomic.Uint32{}
)

// TestCollectAndSend Тест будет работать при условии, что значение reportInterval - минимум, 1 секунда,
// т.к. за это время все метрики должны успеть уйти на сервер
func TestCollectAndSend(t *testing.T) {
	swapFunc := parseAndCleanConfig()
	defer swapFunc()

	server := httptest.NewServer(router(t))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)

	oldEndpoint := config.Configuration.Endpoint
	config.Configuration.Endpoint = fmt.Sprintf("%s:%s", serverUrl.Hostname(), serverUrl.Port())
	defer func() {
		config.Configuration.Endpoint = oldEndpoint
	}()

	c := client.NewClientWithRetry()
	c.Use(decorator.GzipCompress)

	metricManager, err := New(dictionary.MetricNameList[:], data_collector.New(), c)
	require.NoErrorf(t, err, "Ошибка инициализации менеджера метрик")

	ctx, cancel := context.WithTimeout(context.Background(), config.Configuration.ReportInterval*time.Second+1*time.Second)
	defer cancel()
	_ = metricManager.CollectAndSend(ctx)

	select {
	case <-ctx.Done():
	}

	data := metricManager.Pool().Get().(*interimData)
	testMetricList(t, "test collect method", data.metricList)
	testMetricList(t, "test send method", processedMetricList)
}

func testMetricList(t *testing.T, testName string, metrics []dto.Metrics) {
	t.Run(testName, func(t *testing.T) {
		metricCountMap := make(map[string]int)
		zeroValueMetricCount := 0
		allowedMetricTypes := []string{dictionary.GaugeMetricType, dictionary.CounterMetricType}

		for _, metric := range metrics {
			if _, ok := metricCountMap[metric.ID]; !ok {
				metricCountMap[metric.ID] = 0
			}

			metricCountMap[metric.ID]++

			if metricCountMap[metric.ID] > 1 {
				t.Fatalf("metric '%s' passed more than once", metric.ID)
			}

			assert.Containsf(t, allowedMetricTypes, metric.MType, "metric %s is of a type not allowed", metric.ID)
			assert.Containsf(t, dictionary.MetricNameList, metric.ID, "metric %s is of a name not allowed", metric.ID)

			if metric.MType == dictionary.GaugeMetricType {
				assert.GreaterOrEqualf(t, *metric.Value, 0.0, "%s metric value must be non-negative", metric.ID)
				if *metric.Value == 0.0 {
					zeroValueMetricCount++
				}
				continue
			}

			assert.GreaterOrEqualf(t, *metric.Delta, int64(0), "%s metric value must be non-negative", metric.ID)
			if *metric.Delta == 0 {
				zeroValueMetricCount++
			}
		}

		assert.Equalf(t, len(dictionary.MetricNameList), len(processedMetricList), "Passed %d metrics, need %d", len(processedMetricList), len(dictionary.MetricNameList))

		zeroValueMetricsPercent := float64(zeroValueMetricCount) / float64(len(dictionary.MetricNameList))
		assert.LessOrEqualf(t, zeroValueMetricsPercent, errorRate, "Too many metrics with zero values (%.0f%%)", zeroValueMetricsPercent*100)
	})
}

func router(t *testing.T) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", func(_ http.ResponseWriter, r *http.Request) {
			var mList []dto.Metrics

			body, err := decompressRequestBody(r)
			require.NoError(t, err)

			err = json.NewDecoder(bytes.NewReader(body)).Decode(&mList)
			require.NoError(t, err)

			for _, m := range mList {
				testName := nameOfTestByMetricName(m.ID)
				t.Run(testName, func(t *testing.T) {
					if !assert.Equal(t, http.MethodPost, r.Method) {
						return
					}

					assert.Contains(t, r.Header.Get("Content-Type"), "application/json")
					assert.NoError(t, err)

					mu.Lock()
					processedMetricList = append(processedMetricList, m)
					mu.Unlock()
				})
			}
		})
	})

	return r
}

func decompressRequestBody(r *http.Request) ([]byte, error) {
	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		return []byte{}, errors.New("failed to declare gzip reader")
	}

	body, err := io.ReadAll(gz)
	if err != nil {
		return []byte{}, errors.New("failed to decompress data via gzip writer")
	}

	if err = gz.Close(); err != nil {
		return []byte{}, errors.New("failed to close gzip reader")
	}

	return body, nil
}

func nameOfTestByMetricName(MName string) string {
	requestCount.Swap(requestCount.Load() + 1)

	if MName != "" {
		return fmt.Sprintf("request #%d (%s metric)", requestCount.Load(), MName)
	}
	return fmt.Sprintf("request #%d", requestCount.Load())
}

func parseAndCleanConfig() func() {
	oldReportInterval := config.Configuration.ReportInterval
	config.Configuration.ReportInterval = 2

	return func() {
		config.Configuration.ReportInterval = oldReportInterval
	}
}
