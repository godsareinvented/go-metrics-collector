package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"go.uber.org/multierr"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

const (
	getAllMetricsQuery = `
SELECT
    metrics.ID,
    metric_type.metric_type,
    metrics.delta,
    metrics.value
FROM postgres.public.metrics
JOIN postgres.public.metric_type
	ON metrics.metric_type_id = metric_type.ID;`

	getMetricByIDQuery = `
SELECT
    metrics.ID,
    metric_type.metric_type,
    metrics.delta,
    metrics.value
FROM postgres.public.metrics
JOIN postgres.public.metric_type
	ON metrics.metric_type_id = metric_type.ID
WHERE metrics.ID = $1
	AND metric_type.metric_type = $2;`

	saveOrUpdateMetricQuery = `
WITH metric_type_id_cte AS (
    SELECT id FROM postgres.public.metric_type WHERE metric_type = $4
)

INSERT INTO postgres.public.metrics ("id", "delta", "value", "metric_type_id")
VALUES ($1, $2, $3, (SELECT id FROM metric_type_id_cte))
ON CONFLICT (id) DO UPDATE
  SET ID = $1, delta = $2, value = $3, metric_type_id = (SELECT id FROM metric_type_id_cte);`
)

func (s *PostgreSQLStorage) GetAll(ctx context.Context) ([]dto.Metrics, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	queryRows, err := tx.QueryContext(ctx, getAllMetricsQuery)
	if err != nil {
		errRollback := tx.Rollback()
		return nil, multierr.Combine(errRollback, err)
	}

	var metrics []dto.Metrics
	for queryRows.Next() {
		var metric dto.Metrics
		if err = queryRows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			errRollback := tx.Rollback()
			return nil, multierr.Combine(errRollback, err)
		}
		metrics = append(metrics, metric)
	}

	if err = queryRows.Err(); err != nil {
		errRollback := tx.Rollback()
		return nil, multierr.Combine(errRollback, err)
	}

	if err = queryRows.Close(); err != nil {
		errRollback := tx.Rollback()
		return metrics, multierr.Combine(errRollback, err)
	}

	if err = tx.Commit(); err != nil {
		errRollback := tx.Rollback()
		return metrics, multierr.Combine(errRollback, err)
	}

	return metrics, nil
}

func (s *PostgreSQLStorage) Get(ctx context.Context, m dto.Metrics) (dto.Metrics, bool, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return dto.Metrics{}, false, err
	}

	queryRow := tx.QueryRowContext(ctx, getMetricByIDQuery, m.ID, m.MType)

	var metric dto.Metrics
	err = queryRow.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if errors.Is(err, sql.ErrNoRows) {
		if err = tx.Commit(); err != nil {
			return dto.Metrics{}, false, err
		}
		return dto.Metrics{}, false, nil
	}
	if err != nil {
		errRollback := tx.Rollback()
		return dto.Metrics{}, false, multierr.Combine(errRollback, err)
	}

	if err = queryRow.Err(); err != nil {
		errRollback := tx.Rollback()
		return dto.Metrics{}, false, multierr.Combine(errRollback, err)
	}

	if err = tx.Commit(); err != nil {
		errRollback := tx.Rollback()
		return metric, false, multierr.Combine(errRollback, err)
	}

	return metric, true, nil
}

func (s *PostgreSQLStorage) Set(ctx context.Context, m dto.Metrics) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, saveOrUpdateMetricQuery, m.ID, m.Delta, m.Value, m.MType)
	if err != nil {
		errTx := tx.Rollback()
		return multierr.Combine(errTx, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		errTx := tx.Rollback()
		return multierr.Combine(errTx, err)
	}
	if rowsAffected == 0 {
		errTx := tx.Rollback()
		return multierr.Combine(errTx, errors.New(fmt.Sprintf("no rows affected by saveOrUpdateMetricQuery %s", m.ID)))
	}

	if err = tx.Commit(); err != nil {
		errRollback := tx.Rollback()
		return multierr.Combine(errRollback, err)
	}

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

func NewStorage(db *sql.DB) interfaces.Storage {
	return &PostgreSQLStorage{db: db}
}
