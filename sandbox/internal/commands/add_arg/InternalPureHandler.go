package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_arg"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addArgAction.AddArg(sandbox, api.ArgProps{
		Path:              props.Path,
		Command:           entries.Target,
		Name:              entries.Name,
		Start:             entries.Start,
		End:               entries.End,
		Type:              entries.Type,
		Required:          entries.Required,
		Default:           entries.Default,
		Trigger:           entries.Trigger,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Description:       entries.Description,
		Position:          entries.Position,
	})
	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
