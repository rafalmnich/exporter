package exporter

import (
	"github.com/dghubble/sling"
	"github.com/jinzhu/gorm"
	"github.com/msales/pkg/v3/clix"
	"golang.org/x/xerrors"
)

var sl *sling.Sling
var db *gorm.DB

func getSling() *sling.Sling {
	if sl == nil {
		sl = sling.New()
	}

	return sl
}

func getDb(ctx *clix.Context) (*gorm.DB, error) {
	if db == nil {
		gormDB, err := newGorm(ctx)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		db = gormDB
	}

	return db, nil
}

func newGorm(ctx *clix.Context) (*gorm.DB, error) {
	gormDB, err := gorm.Open("postgres", ctx.String(flagDBUri))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return gormDB, nil
}
