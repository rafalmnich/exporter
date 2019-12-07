package sink

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/msales/go-clock"
	"github.com/rafalmnich/exporter/progress"
	"github.com/vbauerster/mpb/v4"
)

// Exporter is exporting data fromm readings to database
type Exporter struct {
	db *gorm.DB

	batch int
}

// NewExporter is Exporter constructor
func NewExporter(db *gorm.DB, batch int) *Exporter {
	return &Exporter{db: db, batch: batch}
}

// Export exports given readings as database records
func (e *Exporter) Export(ctx context.Context, readings []*Reading) error {
	prog := progress.NewProgress()
	bar := progress.AddBar("saving to database", prog, int64(len(readings)))
	batch := make([]*Reading, 0, e.batch)

	for i, reading := range readings {
		batch = append(batch, reading)

		if (i+1)%e.batch == 0 {
			err := e.saveBatch(bar, batch)
			if err != nil {
				return fmt.Errorf("couldn't insert data batch: %v", err)
			}
			batch = make([]*Reading, 0, e.batch)
		}
	}

	if len(batch) > 0 {
		err := e.saveBatch(bar, batch)
		if err != nil {
			return fmt.Errorf("couldn't insert data batch: %v", err)
		}
	}

	return e.updateImported(readings[0].Occurred)
}

func (e *Exporter) saveBatch(bar *mpb.Bar, batch []*Reading) error {
	bar.IncrBy(e.batch)
	sql := `INSERT INTO "iqc"."reading" ("name","type","value","occurred") VALUES (
				UNNEST(ARRAY[` + prepareNameArray(batch) + `]),
				UNNEST(ARRAY[` + prepareTypeArray(batch) + `]),
				UNNEST(ARRAY[` + prepareValueArray(batch) + `]),
				UNNEST(ARRAY[` + prepareOccurredArray(batch) + `])
			)
			ON CONFLICT DO NOTHING`

	errs := e.db.Exec(sql).GetErrors()
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

var defaultGlue = ", "
var stringGlue = "', '"
var timestampGlue = "'::timestamp, '"

func prepareNameArray(readings []*Reading) string {
	nameFunc := func(ack *Reading) string {
		return ack.Name
	}

	return prepareArray(readings, nameFunc, stringGlue)
}

func prepareTypeArray(readings []*Reading) string {
	typeFunc := func(ack *Reading) string {
		return strconv.Itoa(int(ack.Type))
	}

	return prepareArray(readings, typeFunc, defaultGlue)
}

func prepareValueArray(readings []*Reading) string {
	valueFunc := func(ack *Reading) string {
		return strconv.Itoa(ack.Value)
	}

	return prepareArray(readings, valueFunc, defaultGlue)
}

func prepareOccurredArray(readings []*Reading) string {
	occuredFunc := func(ack *Reading) string {
		return ack.Occurred.Format("2006-01-02 15:04:05")
	}

	return prepareArray(readings, occuredFunc, timestampGlue)
}

func (e *Exporter) updateImported(occurred time.Time) error {
	readingMorning := getMorning(occurred)
	todayMorning := getMorning(clock.Now())

	if occurred.After(todayMorning) {
		return nil
	}

	imported := &Import{
		Day: readingMorning,
	}

	return e.db.Save(imported).Error
}

func getMorning(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

type extractValue func(*Reading) string

func prepareArray(readings []*Reading, extract extractValue, glue string) string {
	values := make([]string, 0, 0)

	for _, reading := range readings {
		values = append(values, extract(reading))
	}

	ret := strings.Join(values, glue)

	if glue == stringGlue || glue == timestampGlue {
		ret = "'" + ret + "'"
	}

	if glue == timestampGlue {
		ret += "::timestamp"
	}

	return ret
}
