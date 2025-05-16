package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"github.com/stretchr/testify/assert"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
	type (
		requestData struct {
			url    string
			method string
		}

		want struct {
			code        int
			contentType string
		}
	)

	tests := []struct {
		name        string
		requestData requestData
		want        want
	}{
		{
			name:        "positive test #1: valid request data and method / counter type / valid metric name",
			requestData: requestData{url: "/update/counter/CounterMetric/1", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #2: valid metric name",
			requestData: requestData{url: "/update/counter/Metric12/1", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #3: valid metric name with numbers",
			requestData: requestData{url: "/update/counter/CounterMetric/1", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #4: gauge type, integer value",
			requestData: requestData{url: "/update/gauge/GaugeMetric/5", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #5: gauge type, float value",
			requestData: requestData{url: "/update/gauge/GaugeFloatValueMetric/2.67", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #6: overwrite metric",
			requestData: requestData{url: "/update/gauge/GaugeMetric/10.0", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #7 (boundary values): counter type, zero value",
			requestData: requestData{url: "/update/counter/CounterMetric/0", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #8 (boundary values): counter type, max value",
			requestData: requestData{url: "/update/counter/CounterMetric/" + strconv.Itoa(math.MaxInt64), method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #9 (boundary values): gauge type, zero value",
			requestData: requestData{url: "/update/counter/CounterMetric/0", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "negative test #10 (boundary values): gauge type, negative value",
			requestData: requestData{url: "/update/gauge/GaugeMetric/-1", method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "positive test #11 (boundary values): gauge type, max value",
			requestData: requestData{url: "/update/counter/CounterMetric/" + strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64), method: http.MethodPost},
			want:        want{code: http.StatusOK},
		},
		{
			name:        "negative test #1: get method is not allowed",
			requestData: requestData{url: "/update/gauge/GaugeMetric/5", method: http.MethodGet},
			want:        want{code: http.StatusMethodNotAllowed},
		},
		{
			name:        "negative test #2 (boundary values): counter type, negative value",
			requestData: requestData{url: "/update/counter/CounterMetric/-1", method: http.MethodPost},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "negative test #3 (boundary values): gauge type, zero real value",
			requestData: requestData{url: "/update/counter/CounterMetric/0.0", method: http.MethodPost},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "negative test #4: invalid handler name",
			requestData: requestData{url: "/updater/gauge/GaugeMetric/5", method: http.MethodPost},
			want:        want{code: http.StatusNotFound},
		},
		{
			name:        "negative test #5: invalid metric type",
			requestData: requestData{url: "/update/unknown/UnknownMetric/12", method: http.MethodPost},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "negative test #6: invalid metric name",
			requestData: requestData{url: "/update/counter/invalid-metric-name/34", method: http.MethodPost},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "negative test #7: invalid metric value",
			requestData: requestData{url: "/update/counter/InvalidMetric/invalid", method: http.MethodPost},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "negative test #8: empty metric name",
			requestData: requestData{url: "/update/counter/", method: http.MethodPost},
			want:        want{code: http.StatusNotFound},
		},
		{
			name:        "negative test #9: empty metric value",
			requestData: requestData{url: "/update/counter/InvalidMetric/", method: http.MethodPost},
			want:        want{code: http.StatusNotFound},
		},
	}

	memStorage := mem_storage.NewInstance()
	repository.NewInstance(memStorage)

	router := chi.NewRouter()
	router.Post("/update/{type}/{name}/{value}", UpdateMetric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.requestData.method, tt.requestData.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}
