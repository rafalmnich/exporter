package exporter

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/rafalmnich/exporter/sink"
)

// Application is the App
type Application interface {
	Importer
	Exporter
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
	db       *gorm.DB
}

func (a App) IsHealthy() error {
	return a.db.DB().Ping()
}

// NewApplication creates App
func NewApplication(importer Importer, exporter Exporter, db *gorm.DB) *App {
	return &App{importer: importer, exporter: exporter, db: db}
}

// Import imports data
func (a App) Import(ctx context.Context) ([]*sink.Reading, error) {
	return a.importer.Import(ctx)
}

// Export exports data
func (a App) Export(ctx context.Context, imp []*sink.Reading) error {
	return a.exporter.Export(ctx, imp)
}
