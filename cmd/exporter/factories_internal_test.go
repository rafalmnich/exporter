package exporter

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetSling(t *testing.T) {
	s := getSling()

	assert.Equal(t, s, getSling())
}

var pgurl = "postgres://iqcc_user:iqcc_pass@localhost:13302/iqcc?sslmode=disable"

func TestGetDb(t *testing.T) {
	db = nil
	db, err := getDb(pgurl)
	assert.NoError(t, err)

	got, err := getDb(pgurl)
	assert.NoError(t, err)

	assert.Equal(t, db, got)
}
func TestGetDbErrored(t *testing.T) {
	db = nil

	_, err := getDb("wrong db url")
	assert.Error(t, err)
}

func TestNewGorm(t *testing.T) {
	db = nil

	g, err := newGorm(pgurl)
	assert.NoError(t, err)

	assert.IsType(t, &gorm.DB{}, g)
}
