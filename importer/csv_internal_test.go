package importer

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rafalmnich/exporter/sink"
	"github.com/stretchr/testify/assert"
)

func Test_prepareReading_readError(t *testing.T) {
	c := NewCsvImporter(nil, nil, 0, "")

	resp := &http.Response{
		Body: ioutil.NopCloser(&erroredReaderMock{}),
	}

	_, err := c.prepareReading(context.Background(), resp, sink.Input)
	assert.Error(t, err)
}

func Test_prepareReading_MultiCSV(t *testing.T) {
	c := NewCsvImporter(nil, nil, 0, "")

	wrongCSV := `Data;Hour;In7;In8;In9;
2019-09-20;00:01:24;0;10;0;
2019-09-20;00:02:24;10;0;0;
Data;Hour;In7;In8;In9;In10;
2019-09-20;00:03:24;10;0;0;1;
2019-09-20;00:04:24;10;0;0;2;`

	resp := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(wrongCSV)),
	}

	readings, err := c.prepareReading(context.Background(), resp, sink.Input)
	assert.NoError(t, err)
	expected := []*sink.Reading{
		{
			Name:     "In7",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 1, 24, 0, time.UTC),
		},
		{
			Name:     "In8",
			Type:     0,
			Value:    10,
			Occurred: time.Date(2019, 9, 20, 0, 1, 24, 0, time.UTC),
		},
		{
			Name:     "In9",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 1, 24, 0, time.UTC),
		},

		{
			Name:     "In7",
			Type:     0,
			Value:    10,
			Occurred: time.Date(2019, 9, 20, 0, 2, 24, 0, time.UTC),
		},
		{
			Name:     "In8",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 2, 24, 0, time.UTC),
		},
		{
			Name:     "In9",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 2, 24, 0, time.UTC),
		},

		{
			Name:     "In7",
			Type:     0,
			Value:    10,
			Occurred: time.Date(2019, 9, 20, 0, 3, 24, 0, time.UTC),
		},
		{
			Name:     "In8",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 3, 24, 0, time.UTC),
		},
		{
			Name:     "In9",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 3, 24, 0, time.UTC),
		},
		{
			Name:     "In10",
			Type:     0,
			Value:    1,
			Occurred: time.Date(2019, 9, 20, 0, 3, 24, 0, time.UTC),
		},

		{
			Name:     "In7",
			Type:     0,
			Value:    10,
			Occurred: time.Date(2019, 9, 20, 0, 4, 24, 0, time.UTC),
		},
		{
			Name:     "In8",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 4, 24, 0, time.UTC),
		},
		{
			Name:     "In9",
			Type:     0,
			Value:    0,
			Occurred: time.Date(2019, 9, 20, 0, 4, 24, 0, time.UTC),
		},
		{
			Name:     "In10",
			Type:     0,
			Value:    2,
			Occurred: time.Date(2019, 9, 20, 0, 4, 24, 0, time.UTC),
		},
	}

	assert.Equal(t, expected, readings)
}

type erroredReaderMock struct {
}

func (e *erroredReaderMock) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
