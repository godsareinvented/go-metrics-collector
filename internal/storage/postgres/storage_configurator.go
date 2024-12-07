package postgres

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
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
	driver, err := postgres.WithInstance(c.db, &postgres.Config{})
	if nil != err {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration/postgres/", "metrics", driver)
	if nil != err {
		return err
	}

	if err = m.Up(); nil != err && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
