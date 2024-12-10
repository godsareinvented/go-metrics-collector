package postgressql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

const (
	getAllMetricsQuery = `
SELECT
    metric.ID,
    metric_type.metric_type,
    metric.metric_name,
    metric.delta,
    metric.value
FROM postgres.public.metric
JOIN postgres.public.metric_type
	ON metric_type.ID = metric.ID;`

	getMetricByIDQuery = `
SELECT
    metric.ID,
    metric_type.metric_type,
    metric.metric_name,
    metric.delta,
    metric.value
FROM postgres.public.metric
JOIN postgres.public.metric_type
	ON metric_type.ID = metric.ID
WHERE metric.ID = $1
	AND metric_type.metric_type = $2;`

	getMetricByNameQuery = `
SELECT
    metric.ID,
    metric_type.metric_type,
    metric.metric_name,
    metric.delta,
    metric.value
FROM postgres.public.metric
JOIN postgres.public.metric_type
	ON metric_type.ID = metric.ID
WHERE metric.metric_name = $1
	AND metric_type.metric_type = $2;`

	saveOrUpdateMetricQuery = `
INSERT INTO postgres.public.metric ("id", "metric_name", "delta", "value") 
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
  SET ID = $1, metric_name = $2, delta = $3, value = $4;`

	saveOrUpdateMetricTypeQuery = `
INSERT INTO postgres.public.metric_type ("id", "metric_type")
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
	SET ID = $1, metric_type = $2;`

	updateUUIDIsFreeFlagQuery = `
UPDATE postgres.public.uuid SET is_free = false WHERE UUID = $1;`

	getGeneratedIDQuery = `
SELECT get_metric_uuid();`
)

func (s *PostgreSQLStorage) GetAll(ctx context.Context) ([]dto.Metrics, error) {
	queryRows, err := s.db.QueryContext(ctx, getAllMetricsQuery)
	if nil != err {
		return nil, err
	}

	var metrics []dto.Metrics
	for queryRows.Next() {
		var metric dto.Metrics
		err = queryRows.Scan(&metric.ID, &metric.MType, &metric.MName, &metric.Delta, &metric.Value)
		if nil != err {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	err = queryRows.Err()
	if nil != err {
		return metrics, err
	}

	err = queryRows.Close()
	if nil != err {
		return metrics, err
	}

	return metrics, nil
}

func (s *PostgreSQLStorage) GetByID(ctx context.Context, ID string, mType string) (dto.Metrics, bool, error) {
	queryRow := s.db.QueryRowContext(ctx, getMetricByIDQuery, ID, mType)

	var metric dto.Metrics
	err := queryRow.Scan(&metric.ID, &metric.MType, &metric.MName, &metric.Delta, &metric.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.Metrics{}, false, nil
	}
	if nil != err {
		return dto.Metrics{}, false, err
	}

	err = queryRow.Err()
	if nil != err {
		return dto.Metrics{}, false, err
	}

	return metric, true, nil
}

func (s *PostgreSQLStorage) GetByName(ctx context.Context, mName string, mType string) (dto.Metrics, bool, error) {
	queryRow := s.db.QueryRowContext(ctx, getMetricByNameQuery, mName, mType)

	var metric dto.Metrics
	err := queryRow.Scan(&metric.ID, &metric.MType, &metric.MName, &metric.Delta, &metric.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.Metrics{}, false, nil
	}
	if nil != err {
		return dto.Metrics{}, false, err
	}

	err = queryRow.Err()
	if nil != err {
		return dto.Metrics{}, false, err
	}

	return metric, true, nil
}

func (s *PostgreSQLStorage) Save(ctx context.Context, metric dto.Metrics) (string, error) {
	isMetricIDEmpty := false
	if "" == metric.ID {
		ID, err := s.GetGeneratedID(ctx, metric)
		if nil != err {
			return "", err
		}
		metric.ID = ID
		isMetricIDEmpty = true
	}
	_, err := s.db.ExecContext(ctx, saveOrUpdateMetricQuery, metric.ID, metric.MName, metric.Delta, metric.Value)
	if nil != err {
		return "", err
	}
	_, err = s.db.ExecContext(ctx, saveOrUpdateMetricTypeQuery, metric.ID, metric.MType)
	if nil != err {
		return "", err
	}

	if isMetricIDEmpty {
		_, err = s.db.ExecContext(ctx, updateUUIDIsFreeFlagQuery, metric.ID)
		if nil != err {
			return "", err
		}
	}

	return metric.ID, nil
}

func (s *PostgreSQLStorage) GetGeneratedID(ctx context.Context, metric dto.Metrics) (string, error) {
	if "" != metric.ID {
		return metric.ID, nil
	}

	queryRow := s.db.QueryRowContext(ctx, getGeneratedIDQuery)

	var ID string
	err := queryRow.Scan(&ID)
	if nil != err {
		return "", err
	}

	err = queryRow.Err()
	if nil != err {
		return "", err
	}

	return ID, nil
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
