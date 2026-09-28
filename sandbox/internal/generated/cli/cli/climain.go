package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
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
// stdout, which answers ExitOk. A handler that does neither has declined, so
// the next command runs — that is the whole of what makes a middleware a
// middleware. A handler that returns an error without answering ends the chain
// too, through HandleFailure. Every command of one line is handed one
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
	props := &api.CommandProps{}
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
func answer(sandbox *api.Sandbox, args []string, consumed []bool, props *api.CommandProps, response *api.CommandResponse, status func() (int, bool)) {
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

		// The answer is what ends the chain, so a handler that answered
		// *and* returned something has answered: what it returned is
		// reported and goes no further.
		if _, answered := status(); answered {
			if err != nil {
				sandbox.Deps.Std.Log("command %s: %s \n", bound.Name, err.Error())
			}
			return
		}
		if err != nil {
			cliio.RaiseWithCause(sandbox, bound, api.HandlerFailure, api.ExitFailure, "", "", err.Error())
			return
		}
	}

	failLine(sandbox, args, consumed, props, response, api.NotFoundFailure, api.ExitUsage, "")
}

// failLine raises a failure that belongs to no command — nothing matched the
// line. The handler still gets an api.Command carrying the line and the
// response and nothing else, so every Handle* file reads the same whichever
// failure brought it there.
//
// It carries no message on purpose: the wording is the project's, filled in by
// the Handle* file through cliio.FailureOf, which is what makes editing that
// file change what the cli says.
func failLine(sandbox *api.Sandbox, args []string, consumed []bool, props *api.CommandProps, response *api.CommandResponse, kind api.CommandFailureKind, code int, cause string) {
	command := api.NewCommand()
	command.Argv = args
	command.Consumed = consumed
	command.Props = props
	command.Response = response

	cliio.RaiseWithCause(sandbox, command, kind, code, "", "", cause)
}

// recoverCommand turns a panicking handler into one answered command line, so
// the panic is reported through HandleFailure rather than as a stack trace. A
// handler that panicked after answering has already answered: the panic is
// reported and nothing is written over it.
func recoverCommand(sandbox *api.Sandbox, args []string, consumed []bool, props *api.CommandProps, response *api.CommandResponse, status func() (int, bool)) {
	failure := recover()
	if failure == nil {
		return
	}

	sandbox.Deps.Std.Error("command panicked: %v\n", failure)
	if _, answered := status(); answered {
		return
	}

	failLine(sandbox, args, consumed, props, response, api.HandlerFailure, api.ExitFailure,
		sandbox.Deps.Std.Sprintf("%v", failure))
}
