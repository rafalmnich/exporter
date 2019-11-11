package importer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/jinzhu/gorm"
	"github.com/msales/go-clock"
	"github.com/msales/pkg/v3/log"
	"golang.org/x/xerrors"

	"github.com/rafalmnich/exporter/sink"
)

// CsvImporter is a service for importing data from csv file that is online
type CsvImporter struct {
	db          *gorm.DB
	doer        sling.Doer
	baseUri     string
	startOffset time.Duration
}

// NewCsvImporter is CsvImporter constructor
func NewCsvImporter(db *gorm.DB, doer sling.Doer, startOffset time.Duration, baseUri string) *CsvImporter {
	return &CsvImporter{db: db, doer: doer, startOffset: startOffset, baseUri: baseUri}
}

// Import imports data (inputs and outputs) from given mass
func (c *CsvImporter) Import(ctx context.Context) ([]*sink.Reading, error) {
	return c.getNewReadings(ctx, c.getLastSync(), sink.Input)
}

func (c *CsvImporter) getLastSync() *sink.Import {
	imported := sink.Import{}

	err := c.db.Last(&imported).Error
	if err != nil {
		return nil
	}

	return &imported
}

func (c *CsvImporter) getNewReadings(ctx context.Context, reading *sink.Import, tp sink.Type) ([]*sink.Reading, error) {
	uri := c.baseUri + c.fileName(reading)

	var response *http.Response
	request, err := sling.
		New().
		Get(uri).
		Request()

	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	response, err = c.doer.Do(request)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, xerrors.New("Couldn't read from source: " + uri)
	}

	return c.prepareReading(ctx, response, tp)
}

func (c *CsvImporter) fileName(lastImport *sink.Import) string {
	nextImportDate := c.nextImportDate(lastImport)

	dir := nextImportDate.Format("200601")
	date := nextImportDate.Format("20060102")

	return fmt.Sprintf("/logs/%s/i_%s.csv", dir, date)
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

	if len(records) == 0 {
		return nil, errors.New("empty or wrong reading")
	}

	names := records[0]
	readings := make([]*sink.Reading, 0, len(names)*len(records))

	for rowNumber, row := range records {
		if rowNumber == 0 {
			continue
		}

		rs, err := c.extract(row, ctx, names, tp)
		if err == nil {
			readings = append(readings, rs...)
		}

	}

	return readings, nil
}

func (c *CsvImporter) extract(row []string, ctx context.Context, names []string, tp sink.Type) ([]*sink.Reading, error) {
	dateTime := row[0] + " " + row[1]
	occurred, err := time.Parse("2006-01-02 15:04:05", dateTime)
	if err != nil {
		log.Error(ctx, "Cannot parse time: "+dateTime)
		return nil, xerrors.Errorf(": %w", err)
	}

	readings := make([]*sink.Reading, 0, len(row))

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
	return readings, nil
}

func isDateTimeCell(i int) bool {
	return i < 2
}

func (c *CsvImporter) nextImportDate(lastImport *sink.Import) time.Time {
	if lastImport != nil {
		return lastImport.Day.AddDate(0, 0, 1)
	}

	return clock.Now().Add(-c.startOffset)
}
