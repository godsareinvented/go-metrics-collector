package handler

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowMetricList(t *testing.T) {
	memStorage := mem_storage.NewInstance()
	repository.NewInstance(memStorage)

	router := chi.NewRouter()
	router.Get("/", ShowMetricList)

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

	testHandler(t, "no metrics", []dto.Metric{}, router)

	repos := repository.GetInstance()
	repos.UpdateMetric(dto.Metric{Type: "counter", Name: "PollCount", Delta: 527})
	repos.UpdateMetric(dto.Metric{Type: "gauge", Name: "RandomValue", Value: 0.47})

	metricList := []dto.Metric{
		{Type: "counter", Name: "PollCount", Delta: 527},
		{Type: "gauge", Name: "RandomValue", Value: 0.47},
	}
	testHandler(t, "sorted metric list", metricList, router)
}

func testHandler(t *testing.T, testName string, metrics []dto.Metric, router *chi.Mux) {
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

func htmlBody(metrics []dto.Metric) (string, error) {
	tmpl := template.Must(template.ParseFiles(mainPageTplPath))
	data := struct {
		Items []dto.Metric
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
