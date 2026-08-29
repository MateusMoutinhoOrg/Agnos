package main

import (
	"os"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// 1. Build deps via the standard adapter. It picks the argument vector
	//    for the injected Verb parser — os.Args[1:] — and standard output
	//    for the injected Printf.
	deps := agnosadapter.New("trackerdata")

	// 2. Inject deps into the pure library.
	l := agnoslib.New(deps)

	// 3. Cli.CliMain is the whole command-line interface, and it is just
	//    another field of the lib: it reads the command line through the
	//    injected parser, prints through the injected Printf, and returns
	//    the process exit code. This is the entire body of cmd/main.
	//
	//    Run it the way you would run the installed binary:
	//      go run ./examples/libraryExamples/MainCallSample/MainCallSample.go category add groceries
	//      go run ./examples/libraryExamples/MainCallSample/MainCallSample.go spend groceries "weekly shopping" 84.50
	//      go run ./examples/libraryExamples/MainCallSample/MainCallSample.go balance
	os.Exit(l.Cli.CliMain(os.Args[1:]))
}
