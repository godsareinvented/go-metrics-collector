package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oldhanasong/go-metrics-collector/internal/service/retry"
	"github.com/oldhanasong/go-metrics-collector/internal/service/retry/prepared_option"
)

var (
	pgErr             *pgconn.PgError
	retriableErrorMap = map[string]bool{
		pgerrcode.UniqueViolation: true,
	}
)

func OpenConnection(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("DATABASE_DSN is empty")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func OpenConnectionWithRetry(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("DATABASE_DSN is empty")
	}

	var db *sql.DB
	var err error
	callback := func() (error, bool) {
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			return err, isRetriableError(err)
		}
		return nil, false
	}

	err = retry.DoWithRetry(context.Background(), prepared_option.DefaultFixedDelayListOptions, callback)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func isRetriableError(err error) bool {
	if errors.As(err, &pgErr) {
		return retriableErrorMap[pgErr.Code]
	}
	return false
}
