package errors

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commands/help"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

// HandleNotFound answers every command line no command of the chain answered —
// one that matched nothing at all, and one every matching command declined.
// An empty command line prints the general help, and exits 0.
//
// It is a command handler like any other — it prints through the response and
// sets the exit status itself — and it is **yours**: written once by
// `agnos build` and never regenerated, so whatever you put here is
// what your cli says.
//
// What went wrong is on `command.Failure`, read through cliio.FailureOf so a
// command carrying none still answers something. The command is bound when a
// declared command raised the failure and bare when none did; command.Argv is
// the command line either way.
//
// Answer a failure here; never raise one. cliio.Raise comes back to this file.
func HandleNotFound(sandbox *api.Sandbox, command *api.Command, response *api.CommandResponse) error {
	// An empty command line asks for nothing but the help, so it is
	// answered like `help`: with the screen and exit 0.
	if len(command.Argv) == 0 {
		response.SetStatus(api.ExitOk)
		help.PrintGeneralHelp(sandbox, response)
		return nil
	}

	failure := cliio.FailureOf(command, api.ExitUsage, "")
	response.SetStatus(failure.Status)
	if failure.Message != "" {
		response.Error("%s\n", failure.Message)
		return nil
	}

	name := sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
	response.Error("unknown command %q — run '%s help' to see the available commands\n", command.Argv[0], name)
	return nil
}
