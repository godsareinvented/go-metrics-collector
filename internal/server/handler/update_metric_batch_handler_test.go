package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateMetricBatchHandler(t *testing.T) {
	swapFunc := parseAndCleanConfig()
	defer swapFunc()

	router := chi.NewRouter()
	router.Post("/updates", UpdateMetricBatch(context.Background()))

	t.Run("test bad request", func(t *testing.T) {
		tests := []struct {
			name    string
			reqBody string
		}{
			{name: "empty body", reqBody: ""},
			{name: "invalid json format", reqBody: "{0}"},
			{name: "empty body", reqBody: `[{"id":"test","type":"gauge","delta":0}]`},
		}

		r := httptest.NewRequest(http.MethodPost, "/updates", nil)
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				r.Body = io.NopCloser(strings.NewReader(test.reqBody))
				status, _, _ := sendRequest(router, r)
				assert.Equalf(t, http.StatusBadRequest, status, "Несоответствие статус-кода ответа ожидаемому в обработчике")
			})
		}
	})

	counterValue1 := ptrInt(int64(rand.Int31()))
	counterValue1Doubled := ptrInt(*counterValue1 + *counterValue1)
	counterValue2 := ptrInt(int64(rand.Int31()))
	counterValue3 := ptrInt(int64(rand.Int31()))
	counterValueSummarized23 := ptrInt(*counterValue2 + *counterValue3)
	gaugeValue1 := ptrFloat(rand.Float64())

	counterMetric1 := dto.Metrics{ID: "CounterMetric1", MType: dictionary.CounterMetricType, Delta: counterValue1, Value: nil}
	counterMetric1DoubledValue := dto.Metrics{ID: "CounterMetric1", MType: dictionary.CounterMetricType, Delta: counterValue1Doubled, Value: nil}
	counterMetric2 := dto.Metrics{ID: "CounterMetric2", MType: dictionary.CounterMetricType, Delta: counterValue2, Value: nil}
	counterMetric2AnotherValue := dto.Metrics{ID: "CounterMetric2", MType: dictionary.CounterMetricType, Delta: counterValue3, Value: nil}
	counterMetric2SummarizedValue := dto.Metrics{ID: "CounterMetric2", MType: dictionary.CounterMetricType, Delta: counterValueSummarized23, Value: nil}
	gaugeMetric1 := dto.Metrics{ID: "GaugeMetric1", MType: dictionary.GaugeMetricType, Delta: nil, Value: gaugeValue1}

	t.Run("test success request", func(t *testing.T) {
		tests := []struct {
			name     string
			reqBody  []dto.Metrics
			wantBody []dto.Metrics
		}{
			{
				name:     "update many metrics",
				reqBody:  []dto.Metrics{counterMetric1, gaugeMetric1},
				wantBody: []dto.Metrics{counterMetric1, gaugeMetric1},
			},
			{
				name:     "increase counter metric value",
				reqBody:  []dto.Metrics{counterMetric1},
				wantBody: []dto.Metrics{counterMetric1DoubledValue},
			},
			{
				name:     "update two counter metrics",
				reqBody:  []dto.Metrics{counterMetric2, counterMetric2AnotherValue},
				wantBody: []dto.Metrics{counterMetric2SummarizedValue},
			},
		}

		r := httptest.NewRequest(http.MethodPost, "/updates", nil)
		r.Header.Set("Content-Type", "application/json")

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				b, err := json.Marshal(test.reqBody)
				require.NoErrorf(t, err, "Ошибка кодировки метрик в json-формат")

				r.Body = io.NopCloser(bytes.NewReader(b))
				status, headers, body := sendRequest(router, r)

				respBodyBytes, err := io.ReadAll(body)
				require.NoErrorf(t, err, "Ошибка декодирования тела ответа в срез метрик")

				wantBodyBytes, err := json.Marshal(test.wantBody)
				require.NoErrorf(t, err, "Ошибка кодировки метрик в json-формат")

				assert.Equalf(t, http.StatusOK, status, "Несоответствие статус-кода ответа ожидаемому в обработчике")
				assert.Containsf(t, headers.Get("Content-Type"), "application/json", "Отсутствие заголовка Content-Type: application/json в ответе обработчика")
				assert.JSONEqf(t, string(wantBodyBytes), string(respBodyBytes), "Несоответствие среза метрик, декодированного из ответа, с ожидаемым срезом от обработчика")
			})
		}
	})
}
