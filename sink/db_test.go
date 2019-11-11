package sink_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/msales/go-clock"
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
var reading3 = &sink.Reading{
	Name:     "name3",
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
	e := sink.NewExporter(db, 2)

	_ = clock.Mock(time.Date(2019, 1, 1, 17, 2, 1, 0, time.UTC))
	defer clock.Restore()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ( 
			UNNEST(ARRAY['name1', 'name2']), 
			UNNEST(ARRAY[0, 1]), 
			UNNEST(ARRAY[20, 150]), 
			UNNEST(ARRAY['2019-01-01 03:02:01'::timestamp, '2019-01-01 03:02:01'::timestamp]) 
		) ON CONFLICT DO NOTHING`)).
		WillReturnResult(sqlmock.NewResult(2, 1))

	err := e.Export(context.TODO(), input)
	assert.NoError(t, err)
}

func TestExporter_ExportWithSavingLastImport(t *testing.T) {
	input := []*sink.Reading{
		reading1,
		reading2,
	}

	mock, db := tests.MockGormDB()
	e := sink.NewExporter(db, 2)

	occuredMorning := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	_ = clock.Mock(time.Date(2019, 1, 10, 17, 2, 1, 0, time.UTC))
	defer clock.Restore()

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES (
			UNNEST(ARRAY['name1', 'name2']),
			UNNEST(ARRAY[0, 1]),
			UNNEST(ARRAY[20, 150]),
			UNNEST(ARRAY['2019-01-01 03:02:01'::timestamp, '2019-01-01 03:02:01'::timestamp])
		)
		ON CONFLICT DO NOTHING`)).
		//WithArgs("'name1', 'name2'", "0, 1", "20, 150", "'2019-01-01 03:02:01 +0000 UTC', '2019-01-01 03:02:01 +0000 UTC'").
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "iqc"."import" ("day") VALUES ($1) RETURNING "iqc"."import"."id"`)).
		WithArgs(occuredMorning).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	err := e.Export(context.TODO(), input)
	assert.NoError(t, err)
}

func TestExporter_Export_WithErrorOnSave(t *testing.T) {
	input := []*sink.Reading{
		{
			Name:     "test5",
			Type:     sink.Output,
			Value:    120,
			Occurred: occuredTime,
		},
	}
	mock, db := tests.MockGormDB()

	e := sink.NewExporter(db, 1)

	mock.ExpectBegin()

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES (
			UNNEST(ARRAY[$1]),
			UNNEST(ARRAY[$2]),
			UNNEST(ARRAY[$3]),
			UNNEST(ARRAY[$4])
		)
		ON CONFLICT DO NOTHING`)).
		WithArgs("test5", 1, 120, occuredTime).
		WillReturnError(errors.New("test error"))

	err := e.Export(context.TODO(), input)
	assert.Error(t, err)
}
