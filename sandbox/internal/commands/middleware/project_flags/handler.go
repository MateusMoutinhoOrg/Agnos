package project_flags

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

// Handle runs in front of every command: it hands the project
// directory on through props.Path and, on --quiet, turns off the progress
// channel for the rest of the process — results (Printf) and errors (Error)
// still go out, only the "… started with path …" notices stop. It answers
// nothing, so the command after it runs.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	props.Path = input.Path
	if input.Quiet {
		sandbox.Deps.StdDeps.Logf = func(format string, a ...any) (int, error) {
			return 0, nil
		}
	}
	return nil
}
