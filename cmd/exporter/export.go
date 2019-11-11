package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/health"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"gopkg.in/urfave/cli.v1"

	"github.com/rafalmnich/exporter"
	"github.com/rafalmnich/exporter/importer"
	"github.com/rafalmnich/exporter/sink"
)

func run(c *cli.Context) error {
	ctx, err := clix.NewContext(c)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	go stats.RuntimeFromContext(ctx, stats.DefaultRuntimeInterval)

	s := &http.Client{
		Timeout: 2 * time.Second,
	}

	db, err := getDb(c.String(flagDBUri))
	if err != nil {
		panic(err)
	}

	i := importer.NewCsvImporter(db, s, c.Duration(flagStartOffset), c.String(flagBaseUri))
	e := sink.NewExporter(db, c.Int(flagBatchSize))
	app := exporter.NewApplication(i, e, db)

	if err := app.IsHealthy(); err != nil {
		panic(err)
	}

	go health.StartServer(":"+ctx.String(clix.FlagPort), app)

	errs := make(chan error, 0)
	go func() {
		err := getData(ctx, app, c.Duration(flagImportPeriod), c.Bool(flagImportOnlyOnce))
		if err != nil {
			if err == ExportOnceError {
				err = fmt.Errorf("programm exited: %v", err)
			}
			errs <- err
		}
	}()

	go func() {
		if err := anyError(errs); err != nil {
			log.Fatal(ctx, err.Error())
		}
	}()

	<-clix.WaitForSignals()

	log.Info(ctx, "Task finished!")

	return nil
}

func anyError(errs chan error) error {
	return <-errs
}

var ExportOnceError = errors.New("finishing after first iteration due to app setting in IMPORT_ONLY_ONCE=true")

func getData(ctx context.Context, app exporter.Application, tickTime time.Duration, onlyOnce bool) error {
	ticker := getClock().Ticker(tickTime)

	for {
		select {
		case <-ticker.C:
			i, err := app.Import(ctx)
			if err != nil {
				return err
			}

			_ = stats.Inc(ctx, "reading.import.count", int64(len(i)), 1)

			err = app.Export(ctx, i)
			if err != nil {
				return err
			}
			if onlyOnce {
				return ExportOnceError
			}
		}
	}
}
