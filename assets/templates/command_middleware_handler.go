package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// Handle runs in front of every command line matching
// `{{.Pattern}}` on a lower rung of the chain: it is a middleware. Returning
// nil without answering hands the command line to the next command; answering
// — a status, or a print to stdout — ends the chain here. It is not strict, so
// the tokens it reads are ones the command after it does not have to declare.
//
// Refuse a command line by returning Deps.OpinionatedAgnosCli.Fail:
//
//	return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "profile", "unknown profile")
//
// Hand what you learned to the commands after it through props, the command
// line's commandprops.CommandProps — declare the field in
// sandbox/internal/commandprops/project.go:
//
//	props.Profile = input.Profile
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	return nil
}
