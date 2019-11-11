package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const routinesDirectory = "./migration/routine/"
const registryFile = "./migration/registry.go"

func main() {
	routines, _ := ioutil.ReadDir(routinesDirectory)
	migration, _ := os.Create(registryFile)
	_, _ = migration.Write([]byte("package migration \n\nvar routines = map[string]map[int]string{\n"))
	for _, routine := range routines {
		if routine.IsDir() {
			_, _ = migration.Write([]byte("\t\"" + routine.Name() + "\": {\n"))
			versions, _ := ioutil.ReadDir(routinesDirectory + routine.Name())
			for _, version := range versions {
				if strings.HasSuffix(version.Name(), ".sql") {
					_, _ = migration.Write([]byte("\t\t" + strings.TrimSuffix(version.Name(), ".sql") + " : `"))
					sql, _ := os.Open(routinesDirectory + routine.Name() + `/` + version.Name())
					_, _ = io.Copy(migration, sql)
					_, _ = migration.Write([]byte("`,\n"))
				}
			}
			_, _ = migration.Write([]byte("\t},\n"))
		}
	}
	_, _ = migration.Write([]byte("}\n"))
}
