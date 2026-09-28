package rename_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	renameCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	rename_error := renameCommandAction.RenameCommand(sandbox, api.RenameCommandProps{
		Path:    props.Path,
		Command: entries.Target,
		Name:    entries.Name,
	})
	if rename_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", rename_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
