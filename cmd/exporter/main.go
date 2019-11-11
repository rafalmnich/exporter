package main

import (
	"log"
	"os"

	"github.com/msales/pkg/v3/clix"
	"gopkg.in/urfave/cli.v1"

	_ "github.com/joho/godotenv/autoload"
)

const (
	flagBaseUri        = "input-source"
	flagDBUri          = "db-uri"
	flagImportPeriod   = "import-period"
	flagStartOffset    = "start-offset"
	flagImportOnlyOnce = "import-only-once"
	flagBatchSize      = "batch-size"
)

var flags = clix.Flags{
	cli.StringFlag{
		Name:   flagBaseUri,
		Usage:  "Source host for imported data",
		EnvVar: "BASE_URI",
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
	cli.BoolFlag{
		Name:   flagImportOnlyOnce,
		Usage:  "Import the data only once and die - mostly for testing",
		EnvVar: "IMPORT_ONLY_ONCE",
	},
	cli.IntFlag{
		Name:   flagBatchSize,
		Usage:  "Batch size for committing data in sink",
		EnvVar: "BATCH_SIZE",
	},
}.Merge(clix.CommonFlags, clix.ServerFlags)

// Version is the compiled application version.
var Version = "¯\\_(ツ)_/¯"

var commands = []cli.Command{
	{
		Name:   "export",
		Usage:  "Run the exporter",
		Flags:  flags,
		Action: run,
	},
	{
		Name:  "migrate",
		Usage: "Migrate the database.",
		Subcommands: []cli.Command{
			{
				Name:   "up",
				Flags:  flags,
				Action: runMigrateUp,
			},
			{
				Name:   "down",
				Flags:  flags,
				Action: runMigrateDown,
			},
		},
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
