package exporter

import (
	"log"
	"os"

	"github.com/msales/pkg/v3/clix"
	"gopkg.in/urfave/cli.v1"
)

const (
	flagSourceURI    = "input-source"
	flagDBUri        = "db-uri"
	flagImportPeriod = "import-period"
	flagStartOffset  = "start-offset"
)

var flags = clix.Flags{
	cli.StringFlag{
		Name:   flagSourceURI,
		Usage:  "Source uri for exporting data",
		EnvVar: "SOURCE_URI",
	},
	cli.StringFlag{
		Name:   flagDBUri,
		Usage:  "Postgres url",
		EnvVar: "DB_URI",
	},
	cli.DurationFlag{
		Name:   flagImportPeriod,
		Usage:  "The import period - shouldn't be less than 1 min",
		EnvVar: "IMPORT_PERIOD",
	},
	cli.DurationFlag{
		Name:   flagStartOffset,
		Usage:  "How far from now to start getting readings, if no readings in database",
		EnvVar: "START_OFFSET",
	},
}.Merge(clix.CommonFlags, clix.ServerFlags)

// Version is the compiled application version.
var Version = "¯\\_(ツ)_/¯"

var commands = []cli.Command{
	{
		Name:   "server",
		Usage:  "Run the server",
		Flags:  flags,
		Action: run,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "adsrv"
	app.Flags = clix.ProfilerFlags
	app.Before = clix.RunProfiler
	app.Commands = commands
	app.Version = Version

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
