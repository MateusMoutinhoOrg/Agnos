package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
)

// InternalPureHandler answers `{{.Pattern}}`. Every value the command
// declares in command.yaml is already bound on entries; props is what the
// middlewares in front of it set.
//
// Printing to response answers the command line with api.ExitOk; set another
// status with response.SetStatus. Refuse it by returning cliio.Fail, which is
// answered through the project's own handle_failure.go:
//
//	return cliio.Fail(sandbox, api.ExitFailure, "", "nothing to do")
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	response.Printf("{{.Identifier}} called\n")
	return nil
}
