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
	// 1. Build deps via the standard adapter: a real clock, standard output,
	//    the Verb parser over os.Args[1:], a Keep database on disk, and the
	//    assets compiled into this binary for every line the interface says.
	deps := agnosadapter.New(dataPath())

	// 2. Inject them into the pure library.
	l := agnoslib.New(deps)

	// 3. Run the interface and exit with its return — the same os.Args[1:]
	//    the adapter wired the parser over.
	os.Exit(l.Cli.CliMain(os.Args[1:]))
}

// dataPath returns the directory the tracker persists its records under: the
// AGNOS_DATA override when it is set, otherwise a directory in the user's
// home, falling back to the working directory when the home cannot be
// resolved.
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
