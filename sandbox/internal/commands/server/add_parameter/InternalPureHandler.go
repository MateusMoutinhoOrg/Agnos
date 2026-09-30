package add_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_parameter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := addParameterAction.AddParameter(sandbox, api.RouteParameterProps{
		Path:              props.Path,
		Route:             entries.Route,
		Name:              entries.Name,
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
		Position:          entries.Position,
	})
	if err != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
