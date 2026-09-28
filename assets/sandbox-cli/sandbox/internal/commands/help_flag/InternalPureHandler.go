package help_flag

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commands/help"
)

// helpKey is the spelling this middleware reads, and the one a command
// declaring a flag of its own under it keeps: `add-command --help "..."` is the
// help text of the command being declared, not a request for a screen.
const helpKey = "--help"

// InternalPureHandler answers `<command> --help` with the help screen of the
// command the line is for — the next strict command of the chain that matches
// it — and a bare `--help` with the general help. Without --help, or in front
// of a command that declares --help itself, it hands the line on.
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	if !entries.Help {
		return nil
	}

	next := nextCommand(sandbox, entries.FullCommand)
	if next == nil {
		help.PrintGeneralHelp(sandbox, response)
		return nil
	}
	for _, flag := range next.Flags {
		for _, key := range flag.Keys {
			if key == helpKey {
				return nil
			}
		}
	}

	help.PrintCommandHelp(sandbox, response, next)
	return nil
}

// nextCommand is the first strict command of the chain the command line is
// for, nil when there is none.
func nextCommand(sandbox *api.Sandbox, argv []string) *api.Command {
	for _, declared := range sandbox.Cli.Commands {
		if !declared.Strict {
			continue
		}
		bound := api.BindCommand(declared)
		bound.Argv = argv
		if bound.IsActionable(bound) {
			return declared
		}
	}
	return nil
}
