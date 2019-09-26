package exporter

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
)

func Test_export(t *testing.T) {
	assert.NotPanics(t, func() {
		run(initCliContext(getFlags()))
	})
}

func getFlags() map[string]string {
	return map[string]string{
		flagSourceUri: "source.uri",
		flagDBUri:     "db.uri",
	}
}

// initCliContext initializes clix context to be passed to existing application factories.
func initCliContext(args map[string]string) *cli.Context {
	cliArgs := os.Args[0:1]
	for k, v := range args {
		cliArgs = append(cliArgs, fmt.Sprintf("-%s=%s", k, v))
	}

	var cCtx *cli.Context
	app := cli.NewApp()
	app.Flags = flags
	app.Action = func(c *cli.Context) { cCtx = c }

	err := app.Run(cliArgs)
	if err != nil {
		panic(err)
	}

	return cCtx
}
