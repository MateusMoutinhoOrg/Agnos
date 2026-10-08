package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// InternalPureHandler answers `{{.Pattern}}`. Every value the command
// declares in command.yaml is already bound on entries; props is what the
// middlewares in front of it set.
//
// Printing to response answers the command line with api.ExitOk; set another
// status with response.SetStatus. Refuse it by returning Deps.OpinatedAgnosCli.Fail, which is
// answered through the project's own handle_failure.go:
//
//	return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "nothing to do")
func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	response.Printf("{{.Identifier}} called\n")
	return nil
}
