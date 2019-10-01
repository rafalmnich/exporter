package exporter

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/msales/pkg/v3/clix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
)

func Test_Functional_App(t *testing.T) {
	assert.NotPanics(t, func() {
		go run(initCliContext(getFlags()))
		time.Sleep(300 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)
	})
}

func Test_Functional_App_Panics_Context(t *testing.T) {
	assert.Panics(t, func() {
		flags := getFlags()
		flags[clix.FlagLogLevel] = "not existing"

		ctx := initCliContext(flags)

		run(ctx)
	})
}

func getFlags() map[string]string {
	return map[string]string{
		flagSourceURI:     "source.uri",
		flagDBUri:         "postgres://iqcc_user:iqcc_pass@localhost/iqcc?sslmode=disable",
		clix.FlagLogLevel: "info",
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
