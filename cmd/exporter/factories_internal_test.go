package main

import (
	"context"
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/rafalmnich/exporter"
	"github.com/stretchr/testify/assert"
)

func TestGetSling(t *testing.T) {
	s := getSling()

	assert.Equal(t, s, getSling())
}

var pgurl = "postgres://iqcc_user:iqcc_pass@localhost/iqcc?sslmode=disable"

func TestGetDb(t *testing.T) {
	db, err := getDb(pgurl)
	assert.NoError(t, err)

	got, err := getDb(pgurl)
	assert.NoError(t, err)

	assert.Equal(t, db, got)
}

func TestGetDbErrored(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()
	db = nil

	_, err := getDb("wrong db url")
	assert.Error(t, err)
}

func TestNewGorm(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()
	db = nil

	g, err := newGorm(pgurl)
	assert.NoError(t, err)

	assert.IsType(t, &gorm.DB{}, g)
}

func TestGetClock(t *testing.T) {
	c := getClock()

	assert.Equal(t, c, getClock())
}

func TestPanicOnError(t *testing.T) {
	assert.Panics(t, func() {
		panicOnErr(context.Background(), errors.New("test error"))
	})
}

func TestNewMigrator(t *testing.T) {
	ctx, err := initContext(defaultFlags())
	assert.NoError(t, err)

	migrator, err := newMigrator(ctx)
	assert.NoError(t, err)
	assert.IsType(t, &exporter.Migrator{}, migrator)
}

func TestNewMigratorErroredConnection(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	db = nil
	flags := defaultFlags()
	flags[flagDBUri] = "wrong db uri"

	ctx, err := initContext(flags)
	assert.NoError(t, err)

	_, err = newMigrator(ctx)
	assert.Error(t, err)
}

func TestNewMigratorErrored(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	flags := defaultFlags()
	flags[flagDBUri] = "wrong db uri"

	ctx, err := initContext(flags)
	assert.NoError(t, err)

	_, err = newMigrator(ctx)
	assert.Error(t, err)
}
