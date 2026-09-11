package main

import (
	"os"

	agnosadapter "{{.Module}}/adapters/availables/standard"

	agnoslib "{{.Module}}/sandbox"
)

func main() {

	deps := agnosadapter.New()

	lib := agnoslib.New(&deps)
	argslist := os.Args[1:]
	result := lib.Cli.CliMain(argslist)
	os.Exit(result)
}
