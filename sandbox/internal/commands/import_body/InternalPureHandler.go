package import_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	importBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/import_body"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	import_error := importBodyAction.ImportBody(sandbox, api.RouteBodyImportProps{
		Path:        props.Path,
		Route:       entries.Route,
		Json:        entries.Json,
		File:        entries.File,
		Required:    entries.Required,
		Replace:     entries.Replace,
		InferFormat: entries.InferFormat,
	})
	if import_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", import_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
