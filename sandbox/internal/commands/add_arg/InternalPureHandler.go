package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addArgAction.AddArg(sandbox, api.ArgProps{
		Path:        props.Path,
		Command:     entries.Target,
		Name:        entries.Name,
		Description: entries.Description,
		Type:        entries.Type,
		Default:     entries.Default,
		Required:    entries.Required,
		Array:       entries.Array,
		Position:    entries.Position,
	})
	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
