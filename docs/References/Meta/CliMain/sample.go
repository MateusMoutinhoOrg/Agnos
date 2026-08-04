//go:build ignore

// This file is an illustrative sample, not part of the build.
package main

import (
	"os"
	"path/filepath"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

const (
	// dataDirName is where the records live, under the user's home.
	dataDirName = ".examplelib"
	// dataDirEnv overrides that location, so a script can run against
	// state of its own.
	dataDirEnv = "EXAMPLELIB_DATA"
)

// main is the whole executable: wire, run, exit. No command is branched on
// here and nothing is printed here — that all lives inside the sandbox.
func main() {
	// 1. Build deps through the adapter (the opinionated layer).
	deps := agnosadapter.New(dataPath())

	// 2. Inject them into the pure library.
	l := agnoslib.New(deps)

	// 3. Run the interface and exit with its return — the same os.Args[1:]
	//    the adapter wired the argv parser over.
	os.Exit(l.Sandboxmain(os.Args[1:]))
}

// dataPath resolves where state is kept — an OS-bound choice, so it is made
// out here rather than inside the sandbox.
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
