package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
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
	queryRows, err := s.db.QueryContext(ctx, getAllMetricsQuery)
	if err != nil {
		return nil, err
	}

	var metrics []dto.Metrics
	for queryRows.Next() {
		var metric dto.Metrics
		err = queryRows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	err = queryRows.Err()
	if err != nil {
		return metrics, err
	}

	err = queryRows.Close()
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (s *PostgreSQLStorage) Get(ctx context.Context, m dto.Metrics) (dto.Metrics, bool, error) {
	queryRow := s.db.QueryRowContext(ctx, getMetricByIDQuery, m.ID, m.MType)

	var metric dto.Metrics
	err := queryRow.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.Metrics{}, false, nil
	}
	if err != nil {
		return dto.Metrics{}, false, err
	}

	err = queryRow.Err()
	if err != nil {
		return dto.Metrics{}, false, err
	}

	return metric, true, nil
}

func (s *PostgreSQLStorage) Set(ctx context.Context, m dto.Metrics) error {
	res, err := s.db.ExecContext(ctx, saveOrUpdateMetricQuery, m.ID, m.Delta, m.Value, m.MType)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf("no rows affected by saveOrUpdateMetricQuery %s", m.ID))
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
