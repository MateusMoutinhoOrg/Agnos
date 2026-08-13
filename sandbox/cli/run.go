package cli

// The command-line interface of the tracker, written entirely inside the
// sandbox. It reads the command line through the injected Verb argv parser
// (deps.Deps.VerbLib), takes every word it displays from the injected
// embedded assets (deps.Deps.EmbedDeps), and writes every line through the
// injected formatted writer (deps.Deps.Printf), so the whole program stays
// free of OS-bound and third-party imports — the process only hands it an
// argument vector and exits with the code it returns.
//
// No display text is written here: the usage screen, the version, and each
// message below live in files under /assets/, and this package addresses them
// by path. Changing what the interface says is editing an asset, not editing
// Go.
//
// Like sandbox/lib/store, this package is neither an object package nor
// the entry point: it declares no types and no factories, and is called by
// SandboxmainFactory in sandbox/lib.

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands"
)

// Run is the body of api.Lib.Sandboxmain: it dispatches one command line and
// returns the process exit code. args is only read to detect an empty
// command line; the arguments themselves are drained through the injected
// Verb parser, which the adapter wired over the same vector.
func Run(l *api.Lib, args []string) int {
	if len(args) == 0 {
		l.Deps.Printf("%s", config.Usages)
		return api.ExitUsage
	}

	verb := l.Deps.VerbLib
	if verb.IsPresent(config.HelpFlags) {
		l.Deps.Printf("%s", config.Usages)
		return api.ExitOk
	}
	if verb.IsPresent(config.VersionFlags) {
		return commands.VersionCommand(l)
	}
	// Read the flag before the positional arguments: Verb marks a matched
	// flag used, so draining what is left hands back only the command words.
	quiet := verb.IsPresent(config.QuietFlags)

	command, err := verb.GetNextStringArg()
	if err != nil {
		return commands.UsageError(l, config.NoCommand)
	}

	switch command {
	case "help":
		l.Deps.Printf("%s", config.Usages)
		return api.ExitOk
	case "version":
		return commands.VersionCommand(l)
	case "category":
		return commands.Category(l, quiet)
	case "spend":
		return commands.Record(l, api.Spend, quiet)
	case "received":
		return commands.Record(l, api.Received, quiet)
	case "transactions":
		return commands.Transactions(l)
	case "balance":
		return commands.Balance(l)
	}
	return commands.UsageError(l, config.UnknownCommand, command)
}
