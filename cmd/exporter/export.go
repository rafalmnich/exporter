package exporter

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/health"
	"github.com/msales/pkg/v3/log"
	"github.com/msales/pkg/v3/stats"
	"gopkg.in/urfave/cli.v1"

	"github.com/rafalmnich/exporter"
	"github.com/rafalmnich/exporter/importer"
	"github.com/rafalmnich/exporter/sink"
)

func run(c *cli.Context) {
	ctx, err := clix.NewContext(c)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	go stats.RuntimeFromContext(ctx, stats.DefaultRuntimeInterval)

	s := sling.New()
	s.Base(c.String(flagSourceURI))

	db, err := getDb(c.String(flagDBUri))
	if err != nil {
		panic(err)
	}

	i := importer.NewCsvImporter(db, s)
	e := sink.NewExporter(db)
	app := exporter.NewApplication(i, e, db)

	if err := app.IsHealthy(); err != nil {
		panic(err)
	}

	go health.StartServer(":"+ctx.String(clix.FlagPort), app)

	errs := make(chan error, 1)
	go getData(ctx, app, c.Duration(flagImportPeriod), errs)

	<-clix.WaitForSignals()

	log.Info(ctx, "Task finished!")
}

func getData(ctx context.Context, app exporter.Application, tickTime time.Duration, errs chan error) {
	ticker := getClock().Ticker(tickTime)

	for {
		select {
		case <-ticker.C:
			i, err := app.Import(ctx)
			if err != nil {
				errs <- err
				break
			}

			err = app.Export(ctx, i)
			if err != nil {
				errs <- err
				break
			}

			break
		case <-errs:
			ticker.Stop()
			return
		}
	}
}
