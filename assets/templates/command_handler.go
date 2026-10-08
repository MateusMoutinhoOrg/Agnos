package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// Handle answers `{{.Pattern}}`. Every value the command
// declares in command.yaml is already bound on input; props is what the
// middlewares in front of it set.
//
// Printing to response answers the command line with api.ExitOk; set another
// status with response.SetStatus. Refuse it by returning Deps.OpinionatedAgnosCli.Fail, which is
// answered through the project's own handle_failure.go:
//
//	return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "nothing to do")
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	response.Printf("{{.Identifier}} called\n")
	return nil
}
