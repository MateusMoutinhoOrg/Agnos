package add_parameter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addParameterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_parameter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	err := addParameterAction.AddParameter(sandbox, api.AddParameterProps{
		Path:              props.Path,
		Route:             input.Route,
		Name:              input.Name,
		Type:              input.Type,
		Sources:           input.Source,
		Required:          input.Required,
		Default:           input.Default,
		TriggerType:       input.TriggerType,
		TriggerNegate:     input.TriggerNegate,
		TriggerIgnoreCase: input.TriggerIgnoreCase,
		Trigger:           input.Trigger,
		Description:       input.Description,
		Examples:          input.Example,
		Position:          input.Position,
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
