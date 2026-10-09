package main

import (
	"os"

	agnosadapter "{{.Module}}/adapters/bindings/standard"

	agnoslib "{{.Module}}/sandbox"
)

func main() {

	deps := agnosadapter.New()

	lib := agnoslib.New(&deps)
	argslist := os.Args[1:]
	result := lib.Cli.Main(argslist)
	os.Exit(result)
}
