package project

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// InternalPureHandler runs in front of every command: it hands the project
// directory on through props.Path and, on --quiet, turns off the progress
// channel for the rest of the process — results (Printf) and errors (Error)
// still go out, only the "… started with path …" notices stop. It answers
// nothing, so the command after it runs.
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	props.Path = entries.Path
	if entries.Quiet {
		sandbox.Deps.Std.Log = func(format string, a ...any) (int, error) {
			return 0, nil
		}
	}
	return nil
}
