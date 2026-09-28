package set_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_parameter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := setParameterAction.SetParameter(sandbox, api.RouteParameterEditProps{
		Path:              props.Path,
		Route:             entries.Route,
		Name:              entries.Name,
		Rename:            entries.Rename,
		Type:              entries.Type,
		Fonts:             entries.Font,
		Required:          entries.Required,
		Default:           entries.Default,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Trigger:           entries.Trigger,
		Description:       entries.Description,
		Examples:          entries.Example,
		Clear:             entries.Clear,
	})
	if err != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
