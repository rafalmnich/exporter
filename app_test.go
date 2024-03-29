package exporter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rafalmnich/exporter"
	"github.com/rafalmnich/exporter/mocks"
	"github.com/rafalmnich/exporter/sink"
	"github.com/rafalmnich/exporter/tests"
)

func TestNewApplication(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	_, db := tests.MockGormDB()

	app := exporter.NewApplication(imp, exp, db)
	assert.Implements(t, (*exporter.Application)(nil), app)
}

func TestApp_Import(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	expected := []*sink.Reading{
		{
			Name:     "name",
			Type:     sink.Input,
			Value:    20,
			Occurred: time.Now(),
		},
	}
	ctx := context.Background()

	imp.On("Import", ctx).Return(expected, nil)
	app := exporter.NewApplication(imp, exp, nil)

	data, err := app.Import(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}

func TestApp_Export(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	importData := []*sink.Reading{
		{
			Name:     "name",
			Type:     sink.Input,
			Value:    10,
			Occurred: time.Now(),
		},
	}
	ctx := context.Background()

	exp.On("Export", ctx, importData).Return(nil)
	app := exporter.NewApplication(imp, exp, nil)
	err := app.Export(context.Background(), importData)
	assert.NoError(t, err)
}

func TestApp_IsHealthy(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	_, db := tests.MockGormDB()
	app := exporter.NewApplication(imp, exp, db)

	assert.NoError(t, app.IsHealthy())
}
