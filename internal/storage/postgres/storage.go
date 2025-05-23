package postgres

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func NewInstance(dbDsn string) (interfaces.Storage, error) {
	db, err := openedConnection(dbDsn)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLStorage{
		db: db,
	}, nil
}

func openedConnection(dbDsn string) (*sql.DB, error) {
	if "" == dbDsn {
		return nil, errors.New("dbDsn is empty")
	}

	db, err := sql.Open("pgx", dbDsn)
	if nil != err {
		return nil, err
	}

	return db, nil
}
