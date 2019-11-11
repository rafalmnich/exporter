package sink

import (
	"time"
)

// Type is reading type (input, output)
type Type int

const (
	// Input is for input data
	Input Type = iota
	// Output is for output data
	Output
)

// Reading is the reading model
type Reading struct {
	ID       uint64 `gorm:"primary_key"`
	Name     string
	Type     Type
	Value    int
	Occurred time.Time
}

// TableName returns schema table name
func (r *Reading) TableName() string {
	return "iqc.reading"
}

// Import is the import history
type Import struct {
	ID  uint64 `gorm:"primary_key"`
	Day time.Time
}

func (i *Import) TableName() string {
	return "iqc.import"
}
