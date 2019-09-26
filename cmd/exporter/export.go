package exporter

import (
	"github.com/msales/pkg/v3/clix"
	"gopkg.in/urfave/cli.v1"
)

func run(c *cli.Context) {
	ctx, err := clix.NewContext(c)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()
}
