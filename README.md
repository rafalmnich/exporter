# exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/rafalmnich/exporter)](https://goreportcard.com/report/github.com/rafalmnich/exporter)
[![Build Status](https://travis-ci.org/rafalmnich/exporter.svg?branch=master)](https://travis-ci.org/rafalmnich/exporter)
[![Coverage Status](https://coveralls.io/repos/github/rafalmnich/exporter/badge.svg?branch=master)](https://coveralls.io/github/rafalmnich/exporter?branch=master)
[![GoDoc](https://godoc.org/github.com/rafalmnich/exporter?status.svg)](https://godoc.org/github.com/rafalmnich/exporter)
[![GitHub release](https://img.shields.io/github/release/rafalmnich/exporter.svg)](https://github.com/rafalmnich/exporter/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/rafalmnich/exporter/master/LICENSE)

A pet project for exporting csv data from IQControls Mass controller logs to timescale DB. 


##Enviroments

- **BASE_URI** Host for importing the data - usually mass server ip like `http://10.0.0.150`
- **START_OFFSET** How far from now to start getting readings, if no readings in database `time.Duration`, eg.: `1h`
- **IMPORT_PERIOD** The import period - shouldn't be less than 1 min
- **IMPORT_ONLY_ONCE** Import the data only once and die - mostly for testing
