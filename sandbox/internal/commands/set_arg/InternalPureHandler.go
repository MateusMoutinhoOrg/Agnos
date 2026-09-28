package set_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_arg"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setArgAction.SetArg(sandbox, api.ArgEditProps{
		Path:              props.Path,
		Command:           entries.Target,
		Name:              entries.Name,
		Rename:            entries.Rename,
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
		Clear:             entries.Clear,
	})
	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
