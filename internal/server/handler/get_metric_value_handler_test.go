package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetric(t *testing.T) {
	type (
		requestData struct {
			url    string
			method string
		}

		want struct {
			code        int
			contentType string
			mvalue      string
		}
	)

	tests := []struct {
		name        string
		requestData requestData
		want        want
	}{
		{
			name:        "positive test #1: get counter metric",
			requestData: requestData{url: "/value/counter/PollCount", method: http.MethodPost},
			want:        want{code: http.StatusOK, mvalue: "527"},
		},
		{
			name:        "positive test #2: get gauge metric",
			requestData: requestData{url: "/value/gauge/RandomValue", method: http.MethodPost},
			want:        want{code: http.StatusOK, mvalue: "0.47"},
		},
		{
			name:        "negative test #1: get unknown metric",
			requestData: requestData{url: "/value/counter/UnknownMetric", method: http.MethodPost},
			want:        want{code: http.StatusNotFound, mvalue: ""},
		},
		{
			name:        "negative test #2: empty metric name",
			requestData: requestData{url: "/value/counter/", method: http.MethodPost},
			want:        want{code: http.StatusNotFound, mvalue: ""},
		},
		{
			name:        "negative test #3: empty metric data",
			requestData: requestData{url: "/value/", method: http.MethodPost},
			want:        want{code: http.StatusNotFound, mvalue: ""},
		},
	}

	parseAndCleanConfig()

	router := chi.NewRouter()
	router.Post("/value/{type}/{name}", GetMetric)

	repos := config.Configuration.Repository
	repos.UpdateMetric(dto.Metric{Type: "counter", Name: "PollCount", Delta: 527})
	repos.UpdateMetric(dto.Metric{Type: "gauge", Name: "RandomValue", Value: 0.47})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.requestData.method, tt.requestData.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()
			rawBody, err := io.ReadAll(resp.Body)
			body := string(rawBody)

			require.NoError(t, err)

			assert.Equal(t, tt.want.code, resp.StatusCode)
			if tt.want.mvalue != "" {
				assert.Equal(t, tt.want.mvalue, body)
			}
		})
	}
}
