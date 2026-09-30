package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_page"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addPageAction.AddPage(sandbox, api.PageProps{
		Path:  props.Path,
		Name:  entries.Name,
		Title: entries.Title,
	})

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
