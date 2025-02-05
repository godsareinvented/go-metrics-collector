package postgres

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
)

type PostgreSQLConfigurator struct {
	db *sql.DB
}

func (c *PostgreSQLConfigurator) Configure() error {
	return c.createTables()
}

func NewConfigurator(db *sql.DB) interfaces.StorageConfigurator {
	return &PostgreSQLConfigurator{db: db}
}

func (c *PostgreSQLConfigurator) createTables() error {
	return c.applyMigrations()
}

func (c *PostgreSQLConfigurator) applyMigrations() error {
	driver, err := pgx.WithInstance(c.db, &pgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration/postgres/", "metrics", driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
