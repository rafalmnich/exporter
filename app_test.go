package exporter_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rafalmnich/exporter"
	"github.com/rafalmnich/exporter/mocks"
)

func TestNewApplication(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	app := exporter.NewApplication(imp, exp)
	assert.Implements(t, (*exporter.Application)(nil), app)
}

func TestApplication_Import(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	expected := []exporter.ImportData{
		{
			Type: "input",
			Data: map[string]interface{}{
				"foo":   "bar",
				"one":   1.0,
				"two":   2,
				"three": false,
			},
		},
	}
	ctx := context.Background()

	imp.On("Import", ctx).Return(expected)
	app := exporter.NewApplication(imp, exp)

	data, err := app.Import(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}

func TestApplication_Export(t *testing.T) {
	imp := new(mocks.Importer)
	exp := new(mocks.Exporter)
	importData := []exporter.ImportData{
		{
			Type: "input",
			Data: map[string]interface{}{
				"foo":   "bar",
				"one":   1.0,
				"two":   2,
				"three": false,
			},
		},
	}
	ctx := context.Background()

	exp.On("Export", ctx, importData).Return(importData)
	app := exporter.NewApplication(imp, exp)
	err := app.Export(context.Background(), importData)
	assert.NoError(t, err)
}
