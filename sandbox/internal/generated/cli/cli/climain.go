package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cli/command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

// CliMain is the whole dispatch layer: it runs every command of Cli.Commands
// the command line is for, in the order the collector put them — lowest
// `priority` first — and returns the exit status the line was answered with.
// Whether a command is for the line is its own IsActionable's to say; binding
// and running it is its own CommandHandler's.
//
// More than one command may be for one line, which is what a chain is: each one
// runs in turn until one of them **answers** — sets a status, or prints to
// stdout, which answers ExitOk. A middleware that does neither has declined, so
// the next command runs — that is the whole of what makes a middleware a
// middleware; a strict command that does neither ran silently, and answers
// ExitOk. A handler that returns an error ends the chain too, answered or
// not, through HandleFailure. Every command of one line is handed one
// CommandProps and one Consumed, which is how a middleware hands what it
// learned, and the tokens it read, to the commands after it.
//
// Nothing here prints itself. Every way a line can end without a command
// answering it — nothing matched, a value that will not bind, a token nobody
// read, a panic — is raised through cliio.Raise and answered by one of the
// project's own Handle* files.
// Nothing here is generated per command: every command is one declaration
// built by its own NewCommand and collected by
// sandbox/internal/generated/cli/cli/new.go, so this file is the same in every
// project.
func CliMain(sandbox *api.Sandbox, args []string) int {
	response, status := cliio.Tracked(sandbox)
	props := &commandprops.CommandProps{}
	consumed := make([]bool, len(args))

	answer(sandbox, args, consumed, props, response, status)

	code, answered := status()
	if !answered {
		return api.ExitFailure
	}
	return code
}

// answer runs the chain for one command line and, when no command answered
// it, raises the failure that says so.
func answer(sandbox *api.Sandbox, args []string, consumed []bool, props *commandprops.CommandProps, response *api.CommandResponse, status func() (int, bool)) {
	defer recoverCommand(sandbox, args, consumed, props, response, status)

	for _, declared := range sandbox.Cli.Commands {
		bound := api.BindCommand(declared)
		bound.Argv = args
		bound.Consumed = consumed
		bound.Props = props
		bound.Response = response

		if !bound.IsActionable(bound) {
			continue
		}

		err := bound.CommandHandler(bound)

		// An error ends the chain whether or not the handler answered first:
		// a command that printed part of its output and then failed has
		// failed, and the status HandleFailure sets replaces the ExitOk its
		// print implied.
		if err != nil {
			cliio.RaiseWithCause(sandbox, bound, api.HandlerFailure, api.ExitFailure, "", "", err.Error())
			return
		}
		if _, answered := status(); answered {
			return
		}
		// A strict command is the end of the line, not a middleware: one
		// that ran and printed nothing — it wrote a file, removed one — has
		// carried the line out, so its silence is ExitOk, never a decline
		// that would read as "unknown command".
		if bound.Strict {
			response.SetStatus(api.ExitOk)
			return
		}
	}

	// A command whose verb the line starts with is the one it was meant
	// for: its args did not fit, which is a usage error naming what was
	// wrong, never "unknown command".
	if near := nearCommand(sandbox, args); near != nil {
		bound := api.BindCommand(near)
		bound.Argv = args
		bound.Consumed = consumed
		bound.Props = props
		bound.Response = response
		cliio.Raise(sandbox, bound, api.BadUsageFailure, api.ExitUsage, "", usageProblem(sandbox, bound))
		return
	}

	failLine(sandbox, args, consumed, props, response, api.NotFoundFailure, api.ExitUsage, "")
}

// nearCommand is the strict command whose verb — the longest of its
// Identifiers — the command line's segments start with, nil when none does.
func nearCommand(sandbox *api.Sandbox, args []string) *api.Command {
	segments, _ := command.SplitArgv(sandbox, args)
	var near *api.Command
	longest := 0
	for _, declared := range sandbox.Cli.Commands {
		if !declared.Strict {
			continue
		}
		for _, identifier := range declared.Identifiers {
			words := sandbox.Deps.Stringsdeps.Fields(identifier)
			if len(words) == 0 || len(words) > len(segments) || len(words) <= longest {
				continue
			}
			if sandbox.Deps.Stringsdeps.Join(segments[:len(words)], " ") == identifier {
				near, longest = declared, len(words)
			}
		}
	}
	return near
}

// usageProblem says what keeps a command line from fitting the command it
// names: an arg that will not convert, one missing, or a segment too many.
func usageProblem(sandbox *api.Sandbox, bound *api.Command) string {
	segments, _ := command.SplitArgv(sandbox, bound.Argv)
	usage := sandbox.Deps.Std.Sprintf(" (usage: %s)", bound.Pattern)

	for _, arg := range bound.Args {
		values, found := command.ArgSlice(sandbox, segments, arg)
		if !found {
			continue
		}
		if _, converts := command.ArgValue(sandbox, arg, values); !converts {
			return sandbox.Deps.Std.Sprintf("%q is not a valid %s for <%s>%s",
				sandbox.Deps.Stringsdeps.Join(values, " "), argTypeName(arg.Type), arg.Id, usage)
		}
	}

	reads := 0
	for _, arg := range bound.Args {
		if arg.End < 0 {
			reads = -1
			break
		}
		if arg.End+1 > reads {
			reads = arg.End + 1
		}
	}
	if bound.Segments > 0 {
		reads = bound.Segments
	}

	if reads >= 0 && len(segments) < reads {
		for _, arg := range bound.Args {
			if _, found := command.ArgSlice(sandbox, segments, arg); !found {
				return sandbox.Deps.Std.Sprintf("missing <%s>%s", arg.Id, usage)
			}
		}
	}
	if reads >= 0 && len(segments) > reads {
		return sandbox.Deps.Std.Sprintf("unexpected argument %q%s", segments[reads], usage)
	}
	return "the command line does not fit the command" + usage
}

// argTypeName is how an arg type reads in a message.
func argTypeName(kind api.ArgType) string {
	switch kind {
	case api.IntegerArg:
		return "integer"
	case api.NumberArg:
		return "number"
	case api.UuidArg:
		return "uuid"
	}
	return "string"
}

// failLine raises a failure that belongs to no command — nothing matched the
// line. The handler still gets an api.Command carrying the line and the
// response and nothing else, so every Handle* file reads the same whichever
// failure brought it there.
//
// It carries no message on purpose: the wording is the project's, filled in by
// the Handle* file through cliio.FailureOf, which is what makes editing that
// file change what the cli says.
func failLine(sandbox *api.Sandbox, args []string, consumed []bool, props *commandprops.CommandProps, response *api.CommandResponse, kind api.CommandFailureKind, code int, cause string) {
	command := api.NewCommand()
	command.Argv = args
	command.Consumed = consumed
	command.Props = props
	command.Response = response

	cliio.RaiseWithCause(sandbox, command, kind, code, "", "", cause)
}

// recoverCommand turns a panicking handler into one answered command line, so
// the panic is reported once, through HandleFailure, rather than as a stack
// trace. A handler that panicked after printing has still failed: the status
// HandleFailure sets replaces the ExitOk its print implied.
func recoverCommand(sandbox *api.Sandbox, args []string, consumed []bool, props *commandprops.CommandProps, response *api.CommandResponse, status func() (int, bool)) {
	failure := recover()
	if failure == nil {
		return
	}

	cause := sandbox.Deps.Std.Sprintf("command panicked: %v", failure)
	failLine(sandbox, args, consumed, props, response, api.HandlerFailure, api.ExitFailure, cause)
}
