package main

import (
	"os"
	"path/filepath"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

const (
	// dataDirName is the directory the tracker's records live in, created
	// under the user's home directory so an installed binary tracks one
	// budget from wherever it is run.
	dataDirName = ".agnos"
	// dataDirEnv is the environment variable that overrides where the
	// records live, so a script can run against a budget of its own.
	dataDirEnv = "AGNOS_DATA"
)

// main is the whole executable: it wires an adapter into the library, hands
// the command line to lib.SandBox.Sandboxmain — the interface itself, which lives
// inside the sandbox — and exits with the code it returns. Choosing where the
// records live is an OS-bound decision, so it is made here, outside the
// sandbox, and never inside it.
func main() {
	// 1. Build deps via the standard adapter: a real clock, standard output,
	//    the Verb parser over os.Args[1:], a Keep database on disk, and the
	//    assets compiled into this binary for every line the interface says.
	deps := agnosadapter.New(dataPath())

	// 2. Inject them into the pure library.
	l := agnoslib.New(deps)

	// 3. Run the interface and exit with its return — the same os.Args[1:]
	//    the adapter wired the parser over.
	os.Exit(l.CliMain(os.Args[1:]))
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
