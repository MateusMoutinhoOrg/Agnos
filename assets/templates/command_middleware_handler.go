package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// InternalPureHandler runs in front of every command line matching
// `{{.Pattern}}` on a lower rung of the chain: it is a middleware. Returning
// nil without answering hands the command line to the next command; answering
// — a status, or a print to stdout — ends the chain here. It is not strict, so
// the tokens it reads are ones the command after it does not have to declare.
//
// Refuse a command line by returning cliio.Fail:
//
//	return cliio.Fail(sandbox, api.ExitFailure, "profile", "unknown profile")
//
// Hand what you learned to the commands after it through props, the command
// line's commandprops.CommandProps — declare the field in
// sandbox/internal/commandprops/commandprops.go:
//
//	props.Profile = entries.Profile
func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	return nil
}
