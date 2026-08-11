package main

import (
	"fmt"
	"os"

	verblib "github.com/MateusMoutinhoOrg/Verb/sandbox"
)

// main is the build-script entry point. It parses the process's own command
// line through the Verb argv-parser library — the same library the standard
// adapter uses — and dispatches to one of two commands:
//
//	rename <new-module-path>   — replace the module path in every Go file
//	build  <targets...>        — cross-compile for the listed targets
//
// The optional --provider flag chooses docker or podman; when omitted, the
// script auto-detects whichever is installed. Run from the repository root:
//
//	go run ./build/main build i32 rpm deb linux86 --provider docker
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./build/main <command> [arguments]")
		fmt.Fprintln(os.Stderr, "commands: rename, build")
		os.Exit(1)
	}

	v := verblib.New(os.Args[1:])

	// Read flags before positional arguments so Verb marks them used and
	// GetNextStringArg skips them.
	provider := ""
	p, providerErr := v.GetStringOption([]string{"--provider"}, 0)
	if providerErr == nil {
		provider = p
	}

	command, err := v.GetNextStringArg()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: no command provided")
		os.Exit(1)
	}

	switch command {
	case "rename":
		newPath, err := v.GetNextStringArg()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: rename requires a new module path")
			fmt.Fprintln(os.Stderr, "usage: go run ./build/main rename <new-module-path>")
			os.Exit(1)
		}
		if err := renameModule(newPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "build":
		targets := []string{}
		for {
			target, err := v.GetNextStringArg()
			if err != nil {
				break
			}
			targets = append(targets, target)
		}
		if len(targets) == 0 {
			fmt.Fprintln(os.Stderr, "error: build requires at least one target")
			fmt.Fprintln(os.Stderr, "valid targets: i32, rpm, deb, linux86")
			os.Exit(1)
		}
		if err := buildTargets(targets, provider); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n", command)
		fmt.Fprintln(os.Stderr, "commands: rename, build")
		os.Exit(1)
	}
}
