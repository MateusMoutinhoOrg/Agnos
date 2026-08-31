package main

import (
	"os"
	"path/filepath"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

const (
	dataDirName = ".agnos"

	dataDirEnv = "AGNOS_DATA"
)

func main() {

	deps := agnosadapter.New(dataPath())

	l := agnoslib.New(&deps)
	l.Build()

}

func dataPath() string {
	if override := os.Getenv(dataDirEnv); override != "" {
		return override
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return dataDirName
	}
	return filepath.Join(home, dataDirName)
}
