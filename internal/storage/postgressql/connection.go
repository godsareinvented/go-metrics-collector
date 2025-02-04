package postgressql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/service/retry"
	"github.com/godsareinvented/go-metrics-collector/internal/service/retry/prepared_option"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// GetOpenedConnection unused
func GetOpenedConnection(dsn string) (*sql.DB, error) {
	if "" == dsn {
		return nil, errors.New("DATABASE_DSN is empty")
	}

	db, err := sql.Open("pgx", dsn)
	if nil != err {
		return nil, err
	}

	return db, nil
}

func GetOpenedConnectionWithRetry(dsn string) (*sql.DB, error) {
	if "" == dsn {
		return nil, errors.New("DATABASE_DSN is empty")
	}

	var db *sql.DB
	var err error
	callback := func() (error, bool) {
		db, err = sql.Open("pgx", dsn)
		if nil != err {
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
	var pgErr *pgconn.PgError
	var retriableErrorMap = map[string]bool{
		pgerrcode.UniqueViolation: true,
	}

	if errors.As(err, &pgErr) {
		return retriableErrorMap[pgErr.Code]
	}
	return false
}
