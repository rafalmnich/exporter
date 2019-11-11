package main

import (
	"github.com/msales/pkg/v3/clix"
	"github.com/msales/pkg/v3/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/rafalmnich/exporter"
)

func runMigrateUp(c *cli.Context) {
	ctx, err := clix.NewContext(c)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	m := getMigrator(ctx)

	err = m.MigrateUp()
	if err != nil {
		panicOnErr(ctx, err)
	}

	log.Info(ctx, "UP migrations finished.")
}

func runMigrateDown(c *cli.Context) {
	ctx, err := clix.NewContext(c)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	m := getMigrator(ctx)
	err = m.MigrateDown()
	if err != nil {
		panicOnErr(ctx, err)
	}

	log.Info(ctx, "DOWN migrations finished.")
}

func getMigrator(ctx *clix.Context) *exporter.Migrator {
	m, err := newMigrator(ctx)
	if err != nil {
		panicOnErr(ctx, err)
	}

	return m
}
