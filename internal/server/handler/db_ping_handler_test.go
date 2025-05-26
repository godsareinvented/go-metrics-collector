package handler

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/mock"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDbPing(t *testing.T) {
	testConnectorInterfaceImplementation(t)

	oldRepos := storageSwap(t)
	defer func() {
		config.Configuration.Repository = oldRepos
	}()

	testResponse(t)
}

func testConnectorInterfaceImplementation(t *testing.T) {
	t.Run("connector interface implementation", func(t *testing.T) {
		var psStorage interfaces.Storage = &postgres.PostgreSQLStorage{}
		_, ok := psStorage.(interfaces.StorageConnector)
		require.Truef(t, ok, "Хранилище postgres не имплементирует интерфейс StorageConnector")
	})
}

// testWorkCorrectness Корректно ли тестировать связь с реальной бд? В моём случае, — да.
// Тесты запускаются на той же машине, где лежит приложение.
// Deprecated
func testWorkCorrectness(t *testing.T, err error) {
	t.Run("work correctness", func(t *testing.T) {
		require.Nilf(t, err, "Ошибка при обращении к хранилищу postgres")
	})
}

func testResponse(t *testing.T) {
	t.Run("response", func(t *testing.T) {
		resp := doRequest()
		assert.Equalf(t, http.StatusOK, resp.StatusCode, "Несоответствие статус-кода ответа ожидаемому в хендлере")
	})
}

func storageSwap(t *testing.T) *repository.Repository {
	ctrl := gomock.NewController(t)
	mockStorage := configuredMockStorage(mock.NewMockStorage(ctrl))

	newRepos := repository.NewInstance(mockStorage)
	oldRepos := config.Configuration.Repository

	config.Configuration.Repository = newRepos
	return oldRepos
}

func configuredMockStorage(ms *mock.MockStorage) *mock.MockStorage {
	ms.EXPECT().Ping(gomock.Any()).Return(true, nil)

	return ms
}

func doRequest() *http.Response {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ping", nil)
	DbPing(context.Background()).ServeHTTP(w, r)

	return w.Result()
}
