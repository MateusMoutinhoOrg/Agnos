package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addPageAction.AddPage(sandbox, api.AddPageProps{
		Path:  props.Path,
		Name:  input.Name,
		Title: input.Title,
	})

	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
