package remove_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeDepAction.RemoveDep(sandbox, api.RemoveDepProps{
		Path:         props.Path,
		Dep:          entries.Dep,
		WithAdapters: entries.WithAdapters,
	})

	if remove_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
