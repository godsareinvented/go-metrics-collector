package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

// GetAll todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetAll() ([]dto.Metrics, error) {
	// ...
	return nil, nil
}

// Get todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) Get(_ dto.Metrics) (dto.Metrics, bool, error) {
	// ...
	return dto.Metrics{}, false, nil
}

// Set todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) Set(_ dto.Metrics) error {
	// ...
	return nil
}

func (s *PostgreSQLStorage) Close() error {
	return s.db.Close()
}

func (s *PostgreSQLStorage) Ping(ctx context.Context) (bool, error) {
	if err := s.db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func NewInstance(dbDsn string) interfaces.Storage {
	return &PostgreSQLStorage{
		db: getOpenedConnection(dbDsn),
	}
}

func getOpenedConnection(dbDsn string) *sql.DB {
	if "" == dbDsn {
		panic(errors.New("DATABASE_DSN is empty"))
	}

	db, err := sql.Open("pgx", dbDsn)
	if nil != err {
		panic(err)
	}

	return db
}
