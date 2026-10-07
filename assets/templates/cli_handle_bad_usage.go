package errors

import (
	"{{.Module}}/sandbox/api"
)

// HandleBadUsage answers a command line that matched a command whose values
// would not bind: a required arg or flag missing, a value that will not
// convert, one out of its bounds, its enum or its pattern.
//
// It is a command handler like any other — it prints through the response and
// sets the exit status itself — and it is **yours**: written once by
// `{{.GeneratorName}} build` and never regenerated, so whatever you put here is
// what your cli says.
//
// What went wrong is on `command.Failure`, read through Deps.OpinatedAgnosCli.FailureOf so a
// command carrying none still answers something. The command is bound when a
// declared command raised the failure and bare when none did; command.Argv is
// the command line either way.
//
// Answer a failure here; never raise one: a failure raised from here comes back to this file.
func HandleBadUsage(sandbox *api.Sandbox, command *api.Command, response *api.CommandResponse) error {
	failure := sandbox.Deps.OpinatedAgnosCli.FailureOf(command, api.ExitUsage, "the command line does not fit the command")
	response.SetStatus(failure.Status)

	response.Error("%s\n", failure.Message)
	return nil
}
