package exporter

import (
	"database/sql"

	"github.com/remind101/migrate"

	"github.com/rafalmnich/exporter/migration"
)

import _ "github.com/jinzhu/gorm/dialects/postgres"

//go:generate go run ./migration/routine/generator.go

type Migrator struct {
	db       *sql.DB
	schema   *Schema
	migrator *migrate.Migrator
}

type Schema struct{}

func (s *Schema) migrations() []migrate.Migration {
	return migration.Migrations
}

func NewMigrator(db *sql.DB) *Migrator {
	m := migrate.NewPostgresMigrator(db)
	m.TransactionMode = migrate.SingleTransaction

	return &Migrator{
		db:       db,
		schema:   &Schema{},
		migrator: m,
	}
}

func (m *Migrator) MigrateUp() error {
	return m.migrator.Exec(migrate.Up, m.migrations()...)
}

func (m *Migrator) MigrateDown() error {
	return m.migrator.Exec(migrate.Down, m.migrations()...)
}

func (m *Migrator) migrations() []migrate.Migration {
	return m.schema.migrations()
}

func (m *Migrator) IsHealthy() error {
	if err := m.db.Ping(); err != nil {
		return err
	}

	return nil
}
