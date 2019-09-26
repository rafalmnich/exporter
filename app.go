package exporter

import (
	"context"

	"github.com/rafalmnich/exporter/sink"
)

type Application interface {
	Import(ctx context.Context) ([]*sink.Reading, error)
	Export(ctx context.Context, imp []*sink.Reading) error
}

type Importer interface {
	Import(ctx context.Context) ([]*sink.Reading, error)
}

type Exporter interface {
	Export(ctx context.Context, imp []*sink.Reading) error
}

type application struct {
	importer Importer
	exporter Exporter
}

func NewApplication(importer Importer, exporter Exporter) *application {
	return &application{importer: importer, exporter: exporter}
}

func (a application) Import(ctx context.Context) ([]*sink.Reading, error) {
	return a.importer.Import(ctx)
}

func (a application) Export(ctx context.Context, imp []*sink.Reading) error {
	return a.exporter.Export(ctx, imp)
}
