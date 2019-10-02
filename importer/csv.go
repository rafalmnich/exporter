package importer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/jinzhu/gorm"
	"github.com/msales/pkg/v3/log"
	"golang.org/x/xerrors"

	"github.com/rafalmnich/exporter/sink"
)

// CsvImporter is a service for importing data from csv file that is online
type CsvImporter struct {
	db    *gorm.DB
	sling *sling.Sling
}

// NewCsvImporter is CsvImporter constructor
func NewCsvImporter(db *gorm.DB, sling *sling.Sling) *CsvImporter {
	return &CsvImporter{db: db, sling: sling}
}

// Import imports data (inputs and outputs) from given mass
func (c *CsvImporter) Import(ctx context.Context) ([]*sink.Reading, error) {
	lastSync, err := c.getLastSync(ctx)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return c.getNewReadings(ctx, lastSync, sink.Input)
}

func (c *CsvImporter) getLastSync(ctx context.Context) (*sink.Reading, error) {
	reading := &sink.Reading{}

	err := c.db.Last(reading).Error
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return reading, nil
}

func (c *CsvImporter) getNewReadings(ctx context.Context, reading *sink.Reading, tp sink.Type) ([]*sink.Reading, error) {
	dir := reading.Occurred.Format("200601")
	date := reading.Occurred.Format("20060102")
	filename := fmt.Sprintf("/logs/%s/i_%s.csv", dir, date)

	response, err := c.sling.
		New().
		Get(filename).
		ReceiveSuccess(nil)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return c.prepareReading(ctx, response, tp)
}

func (c *CsvImporter) prepareReading(ctx context.Context, response *http.Response, tp sink.Type) ([]*sink.Reading, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	reader := csv.NewReader(strings.NewReader(string(body)))
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	names := records[0]
	readings := make([]*sink.Reading, 0, len(names)*len(records))

	for rowNumber, row := range records {
		if rowNumber == 0 {
			continue
		}

		dateTime := row[0] + " " + row[1]
		occurred, err := time.Parse("2006-01-02 15:04:05", dateTime)
		if err != nil {
			log.Error(ctx, "Cannot parse time: "+dateTime)
			continue
		}

		for i, value := range row {
			if isDateTimeCell(i) || value == "" {
				continue
			}

			v, err := strconv.Atoi(value)
			if err != nil {
				log.Error(ctx, "Cannot parse value: "+value)
				continue
			}

			reading := &sink.Reading{
				Name:     names[i],
				Type:     tp,
				Value:    v,
				Occurred: occurred,
			}
			readings = append(readings, reading)
		}
	}

	return readings, nil
}

func isDateTimeCell(i int) bool {
	return i < 2
}
