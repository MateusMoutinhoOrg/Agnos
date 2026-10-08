package rename_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	renameCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	rename_error := renameCommandAction.RenameCommand(sandbox, api.RenameCommandProps{
		Path:    props.Path,
		Command: entries.Target,
		Name:    entries.Name,
		// --dir declares no default, so "stay where it is" is no value at
		// all, and the top is spelled "/".
		Dir:    entries.Dir,
		HasDir: entries.Dir != "",
	})
	if rename_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", rename_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
