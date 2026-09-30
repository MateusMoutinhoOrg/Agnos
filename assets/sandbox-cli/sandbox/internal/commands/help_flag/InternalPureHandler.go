package help_flag

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
	"{{.Module}}/sandbox/internal/commands/help"
)

// helpKey is the spelling this middleware reads, and the one a command
// declaring a flag of its own under it keeps: `add-command --help "..."` is the
// help text of the command being declared, not a request for a screen.
const helpKey = "--help"

// InternalPureHandler answers `<command> --help` with the help screen of the
// command the line is for — the next strict command of the chain that matches
// it — and a bare `--help` with the general help. Without --help, or in front
// of a command that declares --help itself and is given a value for it, it
// hands the line on.
func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	if !entries.Help {
		return nil
	}

	next := nextCommand(sandbox, entries.FullCommand)
	if next == nil {
		help.PrintGeneralHelp(sandbox, response)
		return nil
	}
	if declaresHelp(next) && helpHasValue(sandbox, entries.FullCommand) {
		return nil
	}

	help.PrintCommandHelp(sandbox, response, next)
	return nil
}

// declaresHelp reports that command declares a --help flag of its own.
func declaresHelp(command *api.Command) bool {
	for _, flag := range command.Flags {
		for _, key := range flag.Keys {
			if key == helpKey {
				return true
			}
		}
	}
	return false
}

// helpHasValue reports that --help is followed by a value on the command line
// — `add-command x --help "..."` — which is what a command declaring --help
// itself reads. A --help with nothing after it, or another flag, is a request
// for the screen whatever the command declares.
func helpHasValue(sandbox *api.Sandbox, argv []string) bool {
	for index, token := range argv {
		if token != helpKey {
			continue
		}
		return index+1 < len(argv) && !sandbox.Deps.Stringsdeps.HasPrefix(argv[index+1], "-")
	}
	return false
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
