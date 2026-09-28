package errors

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commands/help"
	"{{.Module}}/sandbox/internal/generated/cliio"
)

// HandleNotFound answers every command line no command of the chain answered —
// one that matched nothing at all, and one every matching command declined.
// An empty command line prints the general help.
//
// It is a command handler like any other — it prints through the response and
// sets the exit status itself — and it is **yours**: written once by
// `{{.GeneratorName}} build` and never regenerated, so whatever you put here is
// what your cli says.
//
// What went wrong is on `command.Failure`, read through cliio.FailureOf so a
// command carrying none still answers something. The command is bound when a
// declared command raised the failure and bare when none did; command.Argv is
// the command line either way.
//
// Answer a failure here; never raise one. cliio.Raise comes back to this file.
func HandleNotFound(sandbox *api.Sandbox, command *api.Command, response *api.CommandResponse) error {
	failure := cliio.FailureOf(command, api.ExitUsage, "")
	response.SetStatus(failure.Status)

	if len(command.Argv) == 0 {
		help.PrintGeneralHelp(sandbox, response)
		return nil
	}
	if failure.Message != "" {
		response.Error("%s\n", failure.Message)
		return nil
	}

	name := sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
	response.Error("unknown command %q — run '%s help' to see the available commands\n", command.Argv[0], name)
	return nil
}
