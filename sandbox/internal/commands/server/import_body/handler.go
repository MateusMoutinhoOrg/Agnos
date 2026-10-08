package import_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	importBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/import_body"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	import_error := importBodyAction.ImportBody(sandbox, api.ImportBodyProps{
		Path:        props.Path,
		Route:       input.Name,
		Json:        input.Json,
		File:        input.File,
		Required:    input.Required,
		Replace:     input.Replace,
		InferFormat: input.InferFormat,
	})
	if import_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", import_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
