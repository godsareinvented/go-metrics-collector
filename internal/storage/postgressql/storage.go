package postgressql

import (
	"context"
	"database/sql"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

// GetAll todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetAll(ctx context.Context) ([]dto.Metrics, error) {
	return []dto.Metrics{}, nil
}

// GetByID todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetByID(ctx context.Context, ID string, mType string) (dto.Metrics, bool, error) {
	return dto.Metrics{}, false, nil
}

// GetByName todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetByName(ctx context.Context, mName string, mType string) (dto.Metrics, bool, error) {
	return dto.Metrics{}, false, nil
}

// Save todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) Save(ctx context.Context, metric dto.Metrics) (string, error) {
	return "", nil
}

// GetGeneratedID todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetGeneratedID(ctx context.Context, metric dto.Metrics) string {
	return ""
}

func (s *PostgreSQLStorage) CloseConnect() error {
	return s.db.Close()
}

func (s *PostgreSQLStorage) Ping(ctx context.Context) (bool, error) {
	err := s.db.PingContext(ctx)
	if nil != err {
		return false, err
	}
	return true, nil
}

func NewInstance(db *sql.DB) interfaces.StorageInterface {
	return &PostgreSQLStorage{db: db}
}
