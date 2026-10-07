package errors

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// HandleUnexpectedArg answers a command line carrying a token no command of
// the chain read that is no flag: a word left over once every arg was bound.
//
// It is a command handler like any other — it prints through the response and
// sets the exit status itself — and it is **yours**: written once by
// `agnos build` and never regenerated, so whatever you put here is
// what your cli says.
//
// What went wrong is on `command.Failure`, read through Deps.OpinatedAgnosCli.FailureOf so a
// command carrying none still answers something. The command is bound when a
// declared command raised the failure and bare when none did; command.Argv is
// the command line either way.
//
// Answer a failure here; never raise one: a failure raised from here comes back to this file.
func HandleUnexpectedArg(sandbox *api.Sandbox, command *api.Command, response *api.CommandResponse) error {
	failure := sandbox.Deps.OpinatedAgnosCli.FailureOf(command, api.ExitUsage, "unexpected argument")
	response.SetStatus(failure.Status)

	response.Error("%s\n", failure.Message)
	return nil
}
