package handler

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowMetricList(t *testing.T) {
	oldRepos := parseAndCleanConfig()
	defer func() {
		config.Configuration.Repository = oldRepos
	}()

	router := chi.NewRouter()
	router.Get("/", ShowMetricList(context.Background()))

	oldTpl := mainPageTplPath
	mainPageTplPath = `.\..\..\..\template\main_page.html`
	defer func() {
		mainPageTplPath = oldTpl
	}()

	t.Run("invalid request method", func(t *testing.T) {
		statusCode, _, _, err := sendRequest(router, http.MethodPost)
		require.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	testHandler(t, "no metrics", []dto.Metrics{}, router)

	repos := config.Configuration.Repository
	delta := int64(527)
	value := 0.47
	err := repos.UpdateMetric(dto.Metrics{ID: "PollCount", MType: "counter", Delta: &delta})
	require.NoError(t, err)
	err = repos.UpdateMetric(dto.Metrics{ID: "RandomValue", MType: "gauge", Value: &value})
	require.NoError(t, err)

	metricList := []dto.Metrics{
		{ID: "PollCount", MType: "counter", Delta: &delta},
		{ID: "RandomValue", MType: "gauge", Value: &value},
	}
	testHandler(t, "sorted metric list", metricList, router)
}

func testHandler(t *testing.T, testName string, metrics []dto.Metrics, router *chi.Mux) {
	t.Run(testName, func(t *testing.T) {
		statusCode, contentType, body, err := sendRequest(router, http.MethodGet)
		require.Nil(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Contains(t, contentType, "text/html")

		expectedBody, err := htmlBody(metrics)
		require.Nil(t, err)

		assert.Equal(t, expectedBody, body)
	})
}

func sendRequest(router chi.Router, method string) (int, string, string, error) {
	r := httptest.NewRequest(method, "/", nil)
	r.Header.Set("Accept-Encoding", "")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", "", err
	}

	return resp.StatusCode, resp.Header.Get("Content-Type"), string(rawBody), nil
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
