package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addDocAction.AddDoc(sandbox, api.DocProps{
		Path:        props.Path,
		Name:        entries.Name,
		Description: entries.Description,
		Themes:      entries.Theme,
	})
	if add_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
