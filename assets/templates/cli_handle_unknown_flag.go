package errors

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/generated/cliio"
)

// HandleUnknownFlag answers a command line carrying a token that looks like a
// flag and that no command of the chain read — a typo such as --pathh, which
// would otherwise leave the command running on a default.
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
func HandleUnknownFlag(sandbox *api.Sandbox, command *api.Command, response *api.CommandResponse) error {
	failure := cliio.FailureOf(command, api.ExitUsage, "unknown flag")
	response.SetStatus(failure.Status)

	name := sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
	response.Error("%s — run '%s help' for the accepted flags\n", failure.Message, name)
	return nil
}
