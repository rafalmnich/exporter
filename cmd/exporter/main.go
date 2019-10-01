package exporter

import (
	"log"
	"os"

	"github.com/msales/pkg/v3/clix"
	"gopkg.in/urfave/cli.v1"
)

const (
	flagSourceURI = "input-source"
	flagDBUri     = "db-uri"
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
}.Merge(clix.CommonFlags)

// Version is the compiled application version.
var Version = "¯\\_(ツ)_/¯"

var commands = []cli.Command{
	{
		Name:   "server",
		Usage:  "Run the server",
		Flags:  flags.Merge(clix.CommonFlags, clix.ServerFlags),
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
