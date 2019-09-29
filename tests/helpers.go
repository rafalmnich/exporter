package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

// MockGormDB returns sql mock and gorm db for functional testing
func MockGormDB() (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		panic(err)
	}

	return mock, gormDB
}
