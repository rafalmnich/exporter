package exporter

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"syscall"
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

func Test_Functional_App(t *testing.T) {
	var mock sqlmock.Sqlmock
	mock, db = tests.MockGormDB()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	assert.NotPanics(t, func() {
		go run(initCliContext(getFlags()))
		time.Sleep(10 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)
	})
}

func Test_Functional_App_Errored_Export(t *testing.T) {
	db = nil
	assert.NotPanics(t, func() {
		go run(initCliContext(getFlags()))
		time.Sleep(10 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)
	})
}

func Test_Functional_App_Panics_Context(t *testing.T) {
	db = nil
	assert.Panics(t, func() {
		flags := getFlags()
		flags[clix.FlagLogLevel] = "not existing"

		ctx := initCliContext(flags)

		run(ctx)
	})
}

func Test_Functional_App_Panics_Database(t *testing.T) {
	db = nil
	assert.Panics(t, func() {
		flags := getFlags()
		flags[flagDBUri] = "not existing"

		ctx := initCliContext(flags)

		run(ctx)
	})
}

func Test_Functional_App_Panics_Health(t *testing.T) {
	db, err := getDb(dbURI)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		flags := getFlags()

		ctx := initCliContext(flags)

		_ = db.DB().Close()

		run(ctx)
	})
}

func TestGetData(t *testing.T) {
	now := time.Now()

	ctx, err := initContext(getFlags())
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

	var mock sqlmock.Sqlmock
	mock, db = tests.MockGormDB()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	app := exporter.NewApplication(i, e, db)

	errs := make(chan error, 1)

	go getData(ctx, app, time.Millisecond, errs)

	time.Sleep(10 * time.Millisecond)
	errs <- errors.New("test error")
}

func TestGetData_Errored(t *testing.T) {
	now := time.Now()

	ctx, err := initContext(getFlags())
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

	errs := make(chan error, 1)

	go getData(ctx, app, time.Millisecond, errs)
}

func initContext(flags map[string]string) (context.Context, error) {
	return clix.NewContext(initCliContext(flags))
}

// helpers
var dbURI = "postgres://iqcc_user:iqcc_pass@localhost/iqcc?sslmode=disable"

func getFlags() map[string]string {
	return map[string]string{
		flagSourceURI:     "http://example.com",
		flagDBUri:         dbURI,
		flagImportPeriod:  "5ms",
		clix.FlagLogLevel: "info",
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
