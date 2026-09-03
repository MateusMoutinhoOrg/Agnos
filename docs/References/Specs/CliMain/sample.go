package main

// This file is an illustrative sample, not part of the build. It is what
// assets/cli/cmd/main/main.go renders to in a project whose module is
// github.com/you/my-tool.

import (
	"os"

	agnosadapter "github.com/you/my-tool/adapters/availables/standard"

	agnoslib "github.com/you/my-tool/sandbox"
)

func main() {
	// 1. Build the deps through the standard assembly.
	deps := agnosadapter.New()

	// 2. Inject them into the closed sandbox.
	lib := agnoslib.New(&deps)

	// 3. Run the whole interface over the process's argv; 4. exit with its code.
	os.Exit(lib.Cli.CliMain(os.Args[1:]))
}
