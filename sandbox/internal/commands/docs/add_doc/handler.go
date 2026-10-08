package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDocAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_doc"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addDocAction.AddDoc(sandbox, api.AddDocProps{
		Path:        props.Path,
		Name:        input.Name,
		Description: input.Description,
		Themes:      input.Theme,
	})
	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
