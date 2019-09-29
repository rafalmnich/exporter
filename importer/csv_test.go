package importer_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dghubble/sling"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/rafalmnich/exporter/importer"
	"github.com/rafalmnich/exporter/sink"
	"github.com/rafalmnich/exporter/tests"
)

func TestCsvImporter_Import(t *testing.T) {
	mock, db := tests.MockGormDB()
	response := []byte(`Data;Hour;In7;In8;In9;In10;In11;In12;In13;In14;In15;In31;In32;In33;In41;In42;In43;In51;In52;In53;In61;In62;In63;In81;In82;In83;In91;In92;In93;In245;In246;In247;In248;In253;In254;In255;
2019-09-20;00:01:24;0;10;0;0;6;6;6;6;534;253;0;0;236;3;0;0;0;0;166;127;0;234;23;0;240;113;0;60;180;180;550;0;0;10;
2019-09-20;00:02:24;10;0;0;0;6;6;6;6;510;253;0;0;236;3;0;0;0;0;166;128;0;233;25;0;240;113;0;60;180;180;550;0;0;20;
2019-09-20;00:03:25;20;0;0;0;6;6;6;6;510;253;1;0;236;3;0;0;0;0;166;128;0;233;23;0;239;113;0;60;180;180;550;0;0;30;
2019-09-20;00:04:24;0;0;0;0;4;4;4;4;401;253;2;0;236;3;0;0;0;0;166;128;0;233;23;0;239;113;0;60;180;180;550;0;0;40;
2019-09-20;00:05:24;0;0;0;0;6;6;6;6;508;253;2;0;236;3;0;0;0;0;166;128;0;231;26;0;240;114;0;60;180;180;550;0;0;50;
2019-09-20;00:06:24;0;0;0;0;6;6;6;6;503;253;2;0;236;3;0;0;0;0;166;128;0;231;26;0;240;114;0;60;180;180;550;0;0;60;
2019-09-20;00:07:25;0;0;0;0;6;6;6;6;546;253;2;0;236;3;0;0;0;0;165;126;0;232;23;0;238;112;0;60;180;180;550;0;0;70;
2019-09-20;00:08:24;0;0;0;0;6;6;6;6;548;253;0;0;236;3;0;0;0;0;165;126;0;233;23;0;236;115;0;60;180;180;550;0;0;80;
2019-09-20;00:09:24;0;0;0;0;6;6;6;6;548;253;2;0;236;3;0;0;0;0;166;127;0;233;23;0;236;115;0;60;180;180;550;0;0;90;
2019-09-20;00:10:24;0;0;0;0;6;6;6;6;548;253;2;0;236;4;0;0;0;0;166;126;0;231;23;0;239;114;0;60;180;180;550;0;0;100;
2019-09-20;00:11:25;0;0;0;0;6;6;6;6;548;253;0;0;236;3;0;0;0;0;166;126;0;231;23;0;239;114;0;60;180;180;550;0;0;110;
2019-09-20;00:12:24;0;0;0;0;6;6;6;6;539;253;0;0;236;3;0;0;0;0;166;129;0;231;23;0;238;114;0;60;180;180;550;0;0;120;
2019-09-20;00:13:24;0;0;0;0;6;6;6;6;548;253;1;0;236;3;0;0;0;0;166;129;0;231;23;0;238;113;0;60;180;180;550;0;0;130;
2019-09-20;00:14:24;0;0;0;0;6;6;6;6;548;253;0;0;236;3;0;0;0;0;166;127;0;233;23;0;238;113;0;60;180;180;550;0;0;140;
`)

	sl := mockClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(response), nil
	})
	c := importer.NewCsvImporter(db, sl)
	now := time.Date(2019, 9, 20, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	ctx := context.Background()
	ctx = log.WithLogger(ctx, new(mocks.Logger))
	data, err := c.Import(ctx)
	assert.NoError(t, err)

	expected0 := &sink.Reading{
		Name:     "In7",
		Type:     sink.Input,
		Value:    0,
		Occurred: time.Date(2019, 9, 20, 0, 1, 24, 0, time.UTC),
	}

	expected1 := &sink.Reading{
		Name:     "In8",
		Type:     sink.Input,
		Value:    10,
		Occurred: time.Date(2019, 9, 20, 0, 1, 24, 0, time.UTC),
	}

	expected34 := &sink.Reading{
		Name:     "In7",
		Type:     sink.Input,
		Value:    10,
		Occurred: time.Date(2019, 9, 20, 0, 2, 24, 0, time.UTC),
	}

	assert.Equal(t, expected0, data[0])
	assert.Equal(t, expected1, data[1])
	assert.Equal(t, expected34, data[34])
}

func TestCsvImporter_Import_WithDbError(t *testing.T) {
	mock, db := tests.MockGormDB()
	response := []byte(`Data;Hour;In7;In8;In9;
2019-09-20;00:01:24;0;10;0;`)

	sl := mockClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(response), nil
	})
	c := importer.NewCsvImporter(db, sl)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnError(errors.New("test error"))

	ctx := context.Background()
	ctx = log.WithLogger(ctx, new(mocks.Logger))
	_, err := c.Import(ctx)
	assert.Error(t, err)
}

func TestCsvImporter_Import_WithFetchError(t *testing.T) {
	mock, db := tests.MockGormDB()

	sl := mockClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("test error")
	})
	c := importer.NewCsvImporter(db, sl)
	now := time.Date(2019, 9, 20, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	ctx := context.Background()
	ctx = log.WithLogger(ctx, new(mocks.Logger))
	_, err := c.Import(ctx)
	assert.Error(t, err)
}

func TestCsvImporter_Import_WithResponseError(t *testing.T) {
	mock, db := tests.MockGormDB()

	sl := mockClient(func(req *http.Request) (*http.Response, error) {
		return mockErroredResponse(), nil
	})
	c := importer.NewCsvImporter(db, sl)
	now := time.Date(2019, 9, 20, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	ctx := context.Background()
	ctx = log.WithLogger(ctx, new(mocks.Logger))
	_, err := c.Import(ctx)
	assert.Error(t, err)
}

func TestCsvImporter_Import_WithTimeError(t *testing.T) {
	mock, db := tests.MockGormDB()
	response := []byte(`Data;Hour;In7;In8;In9;
not a date;00:01:24;0;10;0;`)

	sl := mockClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(response), nil
	})
	c := importer.NewCsvImporter(db, sl)
	now := time.Date(2019, 9, 20, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "Cannot parse time: not a date 00:01:24")
	ctx = log.WithLogger(ctx, logger)
	_, err := c.Import(ctx)
	assert.NoError(t, err)
}

func TestCsvImporter_Import_WithValueError(t *testing.T) {
	mock, db := tests.MockGormDB()
	response := []byte(`Data;Hour;In7;In8;In9;
2019-01-01;00:01:24;not an int;10;0;`)

	sl := mockClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(response), nil
	})
	c := importer.NewCsvImporter(db, sl)
	now := time.Date(2019, 9, 20, 10, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "iqc"."reading"  ORDER BY "iqc"."reading"."id" DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "value", "occurred"}).
			AddRow(1, "In81", 0, 210, now))

	ctx := context.Background()
	logger := new(mocks.Logger)
	logger.On("Error", "Cannot parse value: not an int")
	ctx = log.WithLogger(ctx, logger)
	_, err := c.Import(ctx)
	assert.NoError(t, err)
}

// helpers
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func mockClient(fn MockDoer) *sling.Sling {
	return sling.New().Doer(fn)
}

type MockDoer func(req *http.Request) (*http.Response, error)

func (fn MockDoer) Do(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func mockResponse(response []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(response)),
	}
}

func mockErroredResponse() *http.Response {
	closer := ioutil.NopCloser(errReader(0))

	return &http.Response{
		StatusCode: 200,
		Body:       closer,
	}
}
