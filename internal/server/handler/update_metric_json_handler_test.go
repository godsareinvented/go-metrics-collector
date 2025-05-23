package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateMetricHandler(t *testing.T) {
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
			name:        "valid request",
			requestData: requestData{method: http.MethodPost, body: `{"id":"PollCount","type":"counter","delta":1}`},
			want:        want{code: http.StatusOK, body: `{"id":"PollCount","type":"counter","delta":1}`},
		},
		{
			name:        "update gauge type",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue","type":"gauge","value":0.47}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":0.47}`},
		},
		{
			name:        "counter metric zero value",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue","type":"gauge","value":0}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":0}`},
		},
		{
			name:        "counter metric negative value",
			requestData: requestData{method: http.MethodPost, body: `{"id":"PollCount","type":"counter","delta":-2}`},
			want:        want{code: http.StatusOK, body: `{"id":"PollCount","type":"counter","delta":-1}`},
		},
		{
			name:        "gauge metric integer value",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue","type":"gauge","value":12}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":12}`},
		},
		{
			name:        "gauge metric zero value",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue","type":"gauge","value":0}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":0}`},
		},
		{
			name:        "gauge metric negative value",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue","type":"gauge","value":-1}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":-1}`},
		},
		{
			name:        "overwrite counter metric",
			requestData: requestData{method: http.MethodPost, body: `{"id":"PollCount","type":"counter","delta":3}`},
			want:        want{code: http.StatusOK, body: `{"id":"PollCount","type":"counter","delta":2}`},
		},
		{
			name:        "overwrite gauge metric",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue","type":"gauge","value":0.98}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue","type":"gauge","value":0.98}`},
		},
		{
			name:        "update metric not from list",
			requestData: requestData{method: http.MethodPost, body: `{"id":"NewRandomValue","type":"gauge","value":0.98}`},
			want:        want{code: http.StatusOK, body: `{"id":"NewRandomValue","type":"gauge","value":0.98}`},
		},
		{
			name:        "update metric with numbers in name",
			requestData: requestData{method: http.MethodPost, body: `{"id":"RandomValue2","type":"gauge","value":0.98}`},
			want:        want{code: http.StatusOK, body: `{"id":"RandomValue2","type":"gauge","value":0.98}`},
		},
		{
			name:        "update metric with invalid name",
			requestData: requestData{method: http.MethodPost, body: `{"id":"12","type":"gauge","value":0.98}`},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "update metric with invalid name",
			requestData: requestData{method: http.MethodPost, body: `{"id":"+","type":"gauge","value":0.98}`},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "update metric with long name",
			requestData: requestData{method: http.MethodPost, body: `{"id":"VeryLongMetricNameTooLongIamTiredAnotherSeventeenCharacters","type":"gauge","value":0.98}`},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "invalid json format",
			requestData: requestData{method: http.MethodPost, body: `{1}`},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "empty body",
			requestData: requestData{method: http.MethodPost, body: ``},
			want:        want{code: http.StatusBadRequest},
		},
		{
			name:        "not allowed method",
			requestData: requestData{method: http.MethodGet, body: `{"id":"RandomValue","type":"gauge"}`},
			want:        want{code: http.StatusMethodNotAllowed},
		},
		{
			name:        "empty metric data",
			requestData: requestData{method: http.MethodPost, body: `{"id":"","type":"","value":0}`},
			want:        want{code: http.StatusBadRequest},
		},
	}

	oldRepos := parseAndCleanConfig()
	defer func() {
		config.Configuration.Repository = oldRepos
	}()

	router := chi.NewRouter()
	router.Post("/update", UpdateMetricJson(context.Background()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.requestData.body))
			r := httptest.NewRequest(tt.requestData.method, "/update", body)
			r.Header.Set("Accept-Encoding", "")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)

			if tt.want.code != http.StatusOK {
				return
			}

			assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

			var respMetrics, wantMetric dto.Metrics
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&respMetrics))
			assert.NoError(t, json.Unmarshal([]byte(tt.want.body), &wantMetric))
			assert.Equal(t, wantMetric, respMetrics)
		})
	}
}
