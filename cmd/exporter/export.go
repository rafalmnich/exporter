package exporter

import (
	"github.com/dghubble/sling"
	"github.com/msales/pkg/v3/clix"
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
		log.Fatal(ctx, err.Error())
	}

	i := importer.NewCsvImporter(db, s)
	e := sink.NewExporter(db)
	app := exporter.NewApplication(i, e)

	if err := app.IsHealthy(); err != nil {
		log.Fatal(ctx, err.Error())
	}

	getData(app)

	<-clix.WaitForSignals()

	log.Info(ctx, "Task finished!")
}

func getData(app *exporter.App) {
	select {}
}
