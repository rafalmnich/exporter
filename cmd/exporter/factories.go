package main

import (
	"context"
	"database/sql"

	"github.com/benbjohnson/clock"
	"github.com/dghubble/sling"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/juju/errors"
	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/log"
	"golang.org/x/xerrors"

	"github.com/rafalmnich/exporter"
)

var sl *sling.Sling
var db *gorm.DB
var cl clock.Clock

func getSling() *sling.Sling {
	if sl == nil {
		sl = sling.New()
	}

	return sl
}

func getDb(dbURI string) (*gorm.DB, error) {
	if db == nil {
		gormDB, err := newGorm(dbURI)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		db = gormDB
	}

	return db, nil
}

func newGorm(dbURI string) (*gorm.DB, error) {
	gormDB, err := gorm.Open("postgres", dbURI)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return gormDB, nil
}

func getClock() clock.Clock {
	if cl == nil {
		cl = clock.New()
	}

	return cl
}

// Misc =================================
func panicOnErr(ctx context.Context, err error) {
	if err != nil {
		log.Error(ctx, errors.Details(err), "error", true)
		panic(err)
	}
}

// Migrator ===========================
func newMigrator(ctx *clix.Context) (*exporter.Migrator, error) {
	uri := ctx.String(flagDBUri)
	conn, _ := sql.Open("postgres", uri)
	err := conn.Ping()
	if err != nil {
		return nil, err
	}

	return exporter.NewMigrator(conn), nil
}
