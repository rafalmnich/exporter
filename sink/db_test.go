package sink_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/rafalmnich/exporter/sink"
	"github.com/rafalmnich/exporter/tests"
)

var occuredTime = time.Date(2019, 1, 1, 3, 2, 1, 0, time.UTC)

var reading1 = &sink.Reading{
	Name:     "name1",
	Type:     sink.Input,
	Value:    20,
	Occurred: occuredTime,
}
var reading2 = &sink.Reading{
	Name:     "name2",
	Type:     sink.Output,
	Value:    150,
	Occurred: occuredTime,
}

func TestExporter_Export(t *testing.T) {
	input := []*sink.Reading{
		reading1,
		reading2,
	}

	mock, db := tests.MockGormDB()
	e := sink.NewExporter(db)

	mock.ExpectBegin()

	mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1,$2,$3,$4) RETURNING "iqc"."reading"."id"`)).
		WithArgs("name1", 0, 20, occuredTime).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1,$2,$3,$4) RETURNING "iqc"."reading"."id"`)).
		WithArgs("name2", 1, 150, occuredTime).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(2))

	mock.ExpectCommit()
	err := e.Export(context.TODO(), input)
	assert.NoError(t, err)
}

func TestExporter_Export_WithErrorOnCommit(t *testing.T) {
	input := []*sink.Reading{
		reading1,
	}
	mock, db := tests.MockGormDB()

	e := sink.NewExporter(db)

	mock.ExpectBegin()

	mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1,$2,$3,$4) RETURNING "iqc"."reading"."id"`)).
		WithArgs("name1", 0, 20, occuredTime).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(3))

	mock.ExpectCommit().WillReturnError(errors.New("test error"))

	err := e.Export(context.TODO(), input)
	assert.Error(t, err)
}
