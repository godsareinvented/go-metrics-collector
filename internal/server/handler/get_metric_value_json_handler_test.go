package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMetricJson(t *testing.T) {
	type (
		requestData struct {
			method string
			body   string
		}

		want struct {
			code        int
			contentType string
			body        string
		}
	)

	tests := []struct {
		name        string
		requestData requestData
		want        want
	}{
		{
			name:        "positive test #1: get counter metric",
			requestData: requestData{method: http.MethodPost, body: `{"id": "PollCount","type": "counter"}`},
			want:        want{code: http.StatusOK, body: `{"id":"PollCount","type":"counter","delta":527}`},
		},
		{
			name:        "positive test #2: get gauge metric",
			requestData: requestData{method: http.MethodPost, body: `{"id": "RandomValue","type": "gauge"}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":0.47}`},
		},
		{
			name:        "negative test #1: get unknown metric",
			requestData: requestData{method: http.MethodPost, body: `{"id": "Unknown","type": "counter"}`},
			want:        want{code: http.StatusNotFound, body: ``},
		},
		{
			name:        "negative test #2: invalid json format",
			requestData: requestData{method: http.MethodPost, body: `{1}`},
			want:        want{code: http.StatusBadRequest, body: ``},
		},
		{
			name:        "negative test #3: empty body",
			requestData: requestData{method: http.MethodPost, body: ``},
			want:        want{code: http.StatusBadRequest, body: ``},
		},
		{
			name:        "negative test #4: not allowed method",
			requestData: requestData{method: http.MethodGet, body: `{"id": "RandomValue","type": "gauge"}`},
			want:        want{code: http.StatusMethodNotAllowed, body: ``},
		},
		{
			name:        "negative test #5: empty metric data",
			requestData: requestData{method: http.MethodPost, body: `{"id": "","type": ""}`},
			want:        want{code: http.StatusBadRequest, body: ``},
		},
	}

	oldRepos := parseAndCleanConfig()
	defer func() {
		config.Configuration.Repository = oldRepos
	}()

	router := chi.NewRouter()
	router.Post("/value", GetMetricJson(context.Background()))

	err := prepareStorage()
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.requestData.body))
			r := httptest.NewRequest(tt.requestData.method, "/value", body)
			r.Header.Set("Accept-Encoding", "")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			require.NoError(t, err)

			assert.Equal(t, tt.want.code, resp.StatusCode)
			if tt.want.code != http.StatusOK {
				return
			}

			var respMetrics, wantMetric dto.Metrics
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&respMetrics))
			assert.NoError(t, json.Unmarshal([]byte(tt.want.body), &wantMetric))
			assert.Equal(t, wantMetric, respMetrics)
		})
	}
}
