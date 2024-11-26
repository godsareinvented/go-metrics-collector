package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

// GetAll todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetAll() ([]dto.Metrics, error) {
	return []dto.Metrics{}, nil
}

// GetByID todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetByID(ID string, mType string) (dto.Metrics, bool, error) {
	return dto.Metrics{}, false, nil
}

// GetByName todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetByName(mName string, mType string) (dto.Metrics, bool, error) {
	return dto.Metrics{}, false, nil
}

// Save todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) Save(metric dto.Metrics) (string, error) {
	return "", nil
}

// GetGeneratedID todo: заглушка для реализации интерфейса. Позже прописать тело функции.
func (s *PostgreSQLStorage) GetGeneratedID(metric dto.Metrics) string {
	return ""
}

func (s *PostgreSQLStorage) Close() error {
	return s.db.Close()
}

func (s *PostgreSQLStorage) Ping(ctx context.Context) (bool, error) {
	err := s.db.PingContext(ctx)
	if nil != err {
		return false, err
	}
	return true, nil
}

func NewInstance(dbDsn string) interfaces.StorageInterface {
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
