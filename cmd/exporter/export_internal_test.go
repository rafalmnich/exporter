package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/msales/pkg/v3/clix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"

	"github.com/rafalmnich/exporter"
	"github.com/rafalmnich/exporter/mocks"
	"github.com/rafalmnich/exporter/sink"
	"github.com/rafalmnich/exporter/tests"
)

var lock = &sync.Mutex{}

func init() {
	prepareTestServer()
}

func Test_Functional_App(t *testing.T) {
	t.Skip()
	lock.Lock()
	defer lock.Unlock()

	var mock sqlmock.Sqlmock
	mock, db = tests.MockGormDB()

	now := time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))
	mock.ExpectBegin()

	firstResultTime, _ := time.Parse("2006-01-02 15:04:05", "2019-09-20 00:01:24")
	secondResultTime, _ := time.Parse("2006-01-02 15:04:05", "2019-09-20 00:02:24")

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`)).
		WithArgs("In7", 0, 0, firstResultTime).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`)).
		WithArgs("In8", 0, 10, firstResultTime).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`)).
		WithArgs("In7", 0, 10, secondResultTime).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`)).
		WithArgs("In8", 0, 0, secondResultTime).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	mock.MatchExpectationsInOrder(false)

	flags := defaultFlags()
	flags[flagImportOnlyOnce] = "true"
	err := run(initCliContext(flags))
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//time.Sleep(8 * time.Millisecond)
}

func Test_Functional_App_Panics_Context(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()
	db = nil
	assert.Panics(t, func() {
		flags := defaultFlags()
		flags[clix.FlagLogLevel] = "not existing"

		ctx := initCliContext(flags)

		run(ctx)
	})
}

func Test_Functional_App_Panics_Database(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	db = nil
	assert.Panics(t, func() {
		flags := defaultFlags()
		flags[flagDBUri] = "not existing"

		ctx := initCliContext(flags)

		run(ctx)
	})
}

func Test_Functional_App_Panics_Health(t *testing.T) {
	//t.Skip()
	lock.Lock()
	defer lock.Unlock()

	_, db = tests.MockGormDB()

	assert.Panics(t, func() {
		flags := defaultFlags()
		ctx := initCliContext(flags)

		_ = db.DB().Close()

		_ = run(ctx)
		time.Sleep(8 * time.Millisecond)
	})
}

func TestGetData(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	now := time.Now()
	ctx, err := initContext(defaultFlags())
	assert.NoError(t, err)

	i := new(mocks.Importer)
	e := new(mocks.Exporter)

	readings := []*sink.Reading{
		{
			ID:       1,
			Name:     "test1",
			Type:     sink.Input,
			Value:    10,
			Occurred: now,
		},
	}
	i.On("Import", ctx).Return(readings, nil)
	e.On("Export", ctx, readings).Return(nil)

	_, db = tests.MockGormDB()

	app := exporter.NewApplication(i, e, db)

	getData(ctx, app, time.Millisecond, true)
}

func TestGetData_ErroredExporter(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	now := time.Now()
	ctx, err := initContext(defaultFlags())
	assert.NoError(t, err)

	i := new(mocks.Importer)
	e := new(mocks.Exporter)
	readings := []*sink.Reading{
		{
			ID:       1,
			Name:     "test1",
			Type:     sink.Input,
			Value:    10,
			Occurred: now,
		},
	}
	i.On("Import", ctx).Return(readings, nil)
	e.On("Export", ctx, readings).Return(errors.New("test error"))

	var mock sqlmock.Sqlmock
	mock, db = tests.MockGormDB()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	app := exporter.NewApplication(i, e, db)

	err = getData(ctx, app, time.Millisecond, true)
	assert.Error(t, err)
}

func initContext(flags map[string]string) (*clix.Context, error) {
	return clix.NewContext(initCliContext(flags))
}

// helpers
var dbURI = "postgres://iqcc_user:iqcc_pass@localhost/iqcc?sslmode=disable"

func defaultFlags() map[string]string {
	return map[string]string{
		flagBaseUri:        "http://localhost:8088",
		flagDBUri:          dbURI,
		flagImportPeriod:   "5ms",
		flagStartOffset:    "336h",
		clix.FlagLogLevel:  "info",
		flagImportOnlyOnce: "false",
		flagBatchSize:      "10",
	}
}

// initCliContext initializes clix context to be passed to existing application factories.
func initCliContext(args map[string]string) *cli.Context {
	cliArgs := os.Args[0:1]
	for k, v := range args {
		cliArgs = append(cliArgs, fmt.Sprintf("-%s=%s", k, v))
	}

	var cCtx *cli.Context
	app := cli.NewApp()
	app.Flags = flags
	app.Action = func(c *cli.Context) { cCtx = c }

	err := app.Run(cliArgs)
	if err != nil {
		panic(err)
	}

	return cCtx
}

func prepareTestServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := []byte(`Data;Hour;In7;In8;
2019-09-20;00:01:24;0;10;
2019-09-20;00:02:24;10;0;
`)

		_, _ = w.Write(response)
	})
	go func() {
		_ = http.ListenAndServe(":8088", nil)
	}()
}
