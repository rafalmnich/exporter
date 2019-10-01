package exporter

import (
	"github.com/benbjohnson/clock"
	"github.com/dghubble/sling"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
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
