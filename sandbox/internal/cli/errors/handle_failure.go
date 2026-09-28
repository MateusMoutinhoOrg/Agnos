package errors

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

// HandleFailure answers a command line a command could not carry out: a
// handler that refused it through cliio.Fail, one that returned an error
// without answering, and one that panicked.
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
func HandleFailure(sandbox *api.Sandbox, command *api.Command, response *api.CommandResponse) error {
	failure := cliio.FailureOf(command, api.ExitFailure, "")
	response.SetStatus(failure.Status)

	// Cause is what went wrong underneath — an error's text, or the value a
	// handler panicked with. It is printed only when the failure says
	// nothing else.
	if failure.Message == "" {
		failure.Message = failure.Cause
	}
	if failure.Message != "" {
		response.Error("%s\n", failure.Message)
	}
	return nil
}
