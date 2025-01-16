package postgressql

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

type PostgreSQLConfigurator struct {
	Db *sql.DB
}

func (c *PostgreSQLConfigurator) Configure() error {
	return c.createDB()
}

func (c *PostgreSQLConfigurator) createDB() error {
	return c.applyMigrations()
}

func (c *PostgreSQLConfigurator) applyMigrations() error {
	driver, err := postgres.WithInstance(c.Db, &postgres.Config{})
	if nil != err {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/migration/postgresql/", "metric", driver)
	if nil != err {
		return err
	}

	if err = m.Up(); nil != err && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
