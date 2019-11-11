package migration

import "github.com/remind101/migrate"

var Migrations = []migrate.Migration{
	{
		ID: 1,
		Up: migrate.Queries([]string{
			`CREATE SCHEMA IF NOT EXISTS iqc`,

			`CREATE TABLE iqc.reading (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL,
				type VARCHAR(255) NOT NULL,
				value int NOT NULL,
				occurred TIMESTAMP NOT NULL,

				PRIMARY KEY(id)
			)`,
			`CREATE UNIQUE INDEX IDX_reading_unique ON iqc.reading("name", "occurred")`,
		}),
		Down: migrate.Queries([]string{
			`DROP TABLE iqc.reading`,
			`DROP SCHEMA iqc`,
		}),
	},
	{
		ID: 2,
		Up: migrate.Queries([]string{
			`CREATE TABLE iqc.import (
				id BIGSERIAL NOT NULL,
				day date NOT NULL,

				PRIMARY KEY(id)
			)`,
		}),
		Down: migrate.Queries([]string{
			`DROP TABLE iqc.import`,
		}),
	},
}
