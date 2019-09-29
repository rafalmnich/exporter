package sink

import (
	"context"

	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"
)

// Exporter is exporting data fromm readings to database
type Exporter struct {
	db *gorm.DB
}

// NewExporter is Exporter constructor
func NewExporter(db *gorm.DB) *Exporter {
	return &Exporter{db: db}
}

// Export exports given readings as database records
func (e *Exporter) Export(ctx context.Context, readings []*Reading) error {
	t := e.db.Begin()

	for _, reading := range readings {
		if err := t.Save(reading).Error; err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}
	if err := t.Commit().Error; err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
