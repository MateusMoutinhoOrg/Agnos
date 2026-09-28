package cliio

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// Fail is how an InternalPureHandler refuses a command line: it builds the
// failure and returns it as an error, and the handler returns it in turn. A
// handler is handed no command to raise it on, so it is CommandHandler — which
// holds the bound command — that raises what comes back through Raise,
// reaching the project's own handle_failure.go:
//
//	return cliio.Fail(sandbox, api.ExitFailure, "", err.Error())
//
// A message left empty is filled by that file's own wording.
func Fail(sandbox *api.Sandbox, status int, field string, message string) error {
	return FailWithCause(sandbox, status, field, message, "")
}

// FailWithCause is Fail carrying what went wrong underneath — an error's text,
// meant for the log, never for whoever typed the command.
func FailWithCause(sandbox *api.Sandbox, status int, field string, message string, cause string) error {
	return &api.CommandFailure{
		Kind:    api.HandlerFailure,
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// Raise is the one way the cli layer itself raises a failure on a bound
// command — the dispatch, CommandHandler: it records what went wrong on the
// command and hands it to the project's own handler for that kind —
// HandleNotFound, HandleBadUsage, HandleUnknownFlag, HandleUnexpectedArg,
// HandleFailure — through sandbox.Cli.Fail.
//
// It goes through the api rather than calling sandbox/internal/cli/errors
// because the package that dispatches to it, sandbox/internal/generated/cli/cli,
// imports every command package, so no command may import it back. The field
// on the api is what crosses that line, and it is filled by the generated
// sandbox/internal/generated/cli/cli/new.go.
func Raise(sandbox *api.Sandbox, command *api.Command, kind api.CommandFailureKind, status int, field string, message string) error {
	return RaiseWithCause(sandbox, command, kind, status, field, message, "")
}

// RaiseWithCause is Raise carrying what went wrong underneath — an error's
// text, or the value a handler panicked with. The cause reaches the handler on
// command.Failure and is meant for the log.
func RaiseWithCause(sandbox *api.Sandbox, command *api.Command, kind api.CommandFailureKind, status int, field string, message string, cause string) error {
	return RaiseFailure(sandbox, command, &api.CommandFailure{
		Kind:    kind,
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	})
}

// RaiseFailure raises one failure already built — what a handler returned
// through Fail — on the bound command it was returned from.
func RaiseFailure(sandbox *api.Sandbox, command *api.Command, failure *api.CommandFailure) error {
	command.Failure = failure

	// A sandbox whose cli was never built — a command bound by hand rather
	// than by the dispatch — has no handler to reach, so the failure is
	// still reported rather than swallowed.
	if sandbox.Cli.Fail == nil {
		return sandbox.Deps.Std.Errorf("%s", failure.Message)
	}

	return sandbox.Cli.Fail(command)
}

// FailureOf is the failure one of the project's Handle* files is answering. A
// file stands for one kind and one wording, and passes its status and wording
// here as the fallback: a failure that carries its own — a flag that would not
// bind — answers with that, and one that carries none answers with the file's.
// The dispatch raises "no command" with no message at all, so the wording
// comes from handle_not_found.go and changing it there changes what the cli
// says.
func FailureOf(command *api.Command, status int, message string) api.CommandFailure {
	if command.Failure == nil {
		return api.CommandFailure{Status: status, Message: message}
	}

	failure := *command.Failure
	if failure.Status == 0 {
		failure.Status = status
	}
	if failure.Message == "" {
		failure.Message = message
	}
	return failure
}
