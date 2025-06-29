package handler

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowMetricList(t *testing.T) {
	swapFunc := parseAndCleanConfig()
	defer swapFunc()

	router := chi.NewRouter()
	router.Get("/", ShowMetricList(context.Background()))

	oldTpl := mainPageTplPath
	mainPageTplPath = `.\..\..\..\..\template\main_page.html`
	defer func() {
		mainPageTplPath = oldTpl
	}()

	t.Run("invalid request method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set("Accept-Encoding", "")

		statusCode, _, _ := sendRequest(router, r)

		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	testHandler(t, "no metrics", []dto.Metrics{}, router)

	repos := config.Configuration.Repository
	err := repos.UpdateMetric(context.Background(), dto.Metrics{ID: "PollCount", MType: "counter", Delta: ptrInt(527)})
	require.NoError(t, err)
	err = repos.UpdateMetric(context.Background(), dto.Metrics{ID: "RandomValue", MType: "gauge", Value: ptrFloat(0.47)})
	require.NoError(t, err)

	metricList := []dto.Metrics{
		{ID: "PollCount", MType: "counter", Delta: ptrInt(527)},
		{ID: "RandomValue", MType: "gauge", Value: ptrFloat(0.47)},
	}
	testHandler(t, "sorted metric list", metricList, router)
}

func testHandler(t *testing.T, testName string, metrics []dto.Metrics, router *chi.Mux) {
	t.Run(testName, func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Accept-Encoding", "")

		statusCode, headers, body := sendRequest(router, r)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Contains(t, headers.Get("Content-Type"), "text/html")

		expectedBody, err := htmlBody(metrics)
		require.Nil(t, err)

		bodyString, err := io.ReadAll(body)
		require.Nil(t, err)

		assert.Equal(t, expectedBody, string(bodyString))
	})
}

func htmlBody(metrics []dto.Metrics) (string, error) {
	tmpl := template.Must(template.ParseFiles(mainPageTplPath))
	data := struct {
		Items []dto.Metrics
	}{
		Items: metrics,
	}

	buf := new(bytes.Buffer)

	err := tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
