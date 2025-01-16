package postgressql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"go.uber.org/multierr"
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
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})
	if nil != err {
		return "", err
	}

	isMetricIDEmpty := false
	if "" == metric.ID {
		ID, err := getGeneratedIDInTx(ctx, tx, metric)
		if nil != err {
			err = rollbackTransaction(tx, err)
			return "", err
		}
		metric.ID = ID
		isMetricIDEmpty = true
	}
	_, err = tx.ExecContext(ctx, saveOrUpdateMetricQuery, metric.ID, metric.MName, metric.Delta, metric.Value)
	if nil != err {
		err = rollbackTransaction(tx, err)
		return "", err
	}
	_, err = tx.ExecContext(ctx, saveOrUpdateMetricTypeQuery, metric.ID, metric.MType)
	if nil != err {
		err = rollbackTransaction(tx, err)
		return "", err
	}

	if isMetricIDEmpty {
		_, err = tx.ExecContext(ctx, updateUUIDIsFreeFlagQuery, metric.ID)
		if nil != err {
			err = rollbackTransaction(tx, err)
			return "", err
		}
	}

	err = tx.Commit()
	if nil != err {
		err = rollbackTransaction(tx, err) // todo: При ошибочном коммите надо делать роллбек? Звучит комично.
		return "", err
	}

	return metric.ID, nil
}

func (s *PostgreSQLStorage) SaveBatch(ctx context.Context, metricBatch []dto.Metrics) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if nil != err {
		return err
	}
	defer tx.Rollback()

	saveOrUpdateMetricQueryStmt, err := tx.PrepareContext(ctx, saveOrUpdateMetricQuery)
	defer saveOrUpdateMetricQueryStmt.Close()
	if nil != err {
		err = rollbackTransaction(tx, err)
		return err
	}
	saveOrUpdateMetricTypeQueryStmt, err := tx.PrepareContext(ctx, saveOrUpdateMetricTypeQuery)
	defer saveOrUpdateMetricTypeQueryStmt.Close()
	if nil != err {
		err = rollbackTransaction(tx, err)
		return err
	}
	updateUUIDIsFreeFlagQueryStmt, err := tx.PrepareContext(ctx, updateUUIDIsFreeFlagQuery)
	defer updateUUIDIsFreeFlagQueryStmt.Close()
	if nil != err {
		err = rollbackTransaction(tx, err)
		return err
	}

	var isMetricIDEmpty bool
	for _, metric := range metricBatch {
		isMetricIDEmpty = false
		if "" == metric.ID {
			ID, err := getGeneratedIDInTx(ctx, tx, metric)
			if nil != err {
				err = rollbackTransaction(tx, err)
				return err
			}
			metric.ID = ID
			isMetricIDEmpty = true
		}

		_, err = saveOrUpdateMetricQueryStmt.ExecContext(ctx, metric.ID, metric.MName, metric.Delta, metric.Value)
		if nil != err {
			err = rollbackTransaction(tx, err)
			return err
		}
		_, err = saveOrUpdateMetricTypeQueryStmt.ExecContext(ctx, metric.ID, metric.MType)
		if nil != err {
			err = rollbackTransaction(tx, err)
			return err
		}

		if isMetricIDEmpty {
			_, err = updateUUIDIsFreeFlagQueryStmt.ExecContext(ctx, metric.ID)
			if nil != err {
				err = rollbackTransaction(tx, err)
				return err
			}
		}
	}

	err = tx.Commit()
	if nil != err {
		err = rollbackTransaction(tx, err) // todo: При ошибочном коммите надо делать роллбек? Звучит комично.
		return err
	}

	return nil
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

func getGeneratedIDInTx(ctx context.Context, tx *sql.Tx, metric dto.Metrics) (string, error) {
	if "" != metric.ID {
		return metric.ID, nil
	}

	queryRow := tx.QueryRowContext(ctx, getGeneratedIDQuery)

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

func rollbackTransaction(tx *sql.Tx, err error) error {
	var errs error
	multierr.AppendInto(&errs, err)
	err = tx.Rollback()
	if nil != err {
		multierr.AppendInto(&errs, err)
	}
	return errs
}

func NewInstance(db *sql.DB) interfaces.StorageInterface {
	return &PostgreSQLStorage{db: db}
}
