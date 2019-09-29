package exporter

import (
	"context"

	"github.com/rafalmnich/exporter/sink"
)

// Application is the App
type Application interface {
	Import(ctx context.Context) ([]*sink.Reading, error)
	Export(ctx context.Context, imp []*sink.Reading) error
}

// Importer is a data importer
type Importer interface {
	Import(ctx context.Context) ([]*sink.Reading, error)
}

// Exporter is a data exporter
type Exporter interface {
	Export(ctx context.Context, imp []*sink.Reading) error
}

type App struct {
	importer Importer
	exporter Exporter
}

func (a App) IsHealthy() error {
	return nil
}

// NewApplication creates App
func NewApplication(importer Importer, exporter Exporter) *App {
	return &App{importer: importer, exporter: exporter}
}

// Import imports data
func (a App) Import(ctx context.Context) ([]*sink.Reading, error) {
	return a.importer.Import(ctx)
}

// Export exports data
func (a App) Export(ctx context.Context, imp []*sink.Reading) error {
	return a.exporter.Export(ctx, imp)
}
