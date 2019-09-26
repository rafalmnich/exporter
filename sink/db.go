package sink

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Type int

const (
	Input Type = iota
	Output
)

type Reading struct {
	gorm.Model

	ID       int64
	Name     string
	Type     Type
	Value    int
	Occurred time.Time
}
