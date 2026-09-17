package import_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	importBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/import_body"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	import_error := importBodyAction.ImportBody(sandbox, api.RouteBodyImportProps{
		Path:        command.GetString("path"),
		Route:       command.GetString("route"),
		Json:        command.GetString("json"),
		File:        command.GetString("file"),
		Required:    command.GetBool("required"),
		Replace:     command.GetBool("replace"),
		InferFormat: command.GetBool("infer-format"),
	})
	if import_error != nil {
		sandbox.Deps.Std.Error("%s\n", import_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
