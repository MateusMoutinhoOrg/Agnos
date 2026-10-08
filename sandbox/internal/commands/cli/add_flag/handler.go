package add_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addFlagAction.AddFlag(sandbox, api.AddFlagProps{
		Path:              props.Path,
		Command:           input.Command,
		Name:              input.Name,
		Keys:              input.Key,
		Type:              input.Type,
		Required:          input.Required,
		Default:           input.Default,
		Min:               input.Min,
		Max:               input.Max,
		Enum:              input.Enum,
		Pattern:           input.Pattern,
		Trigger:           input.Trigger,
		TriggerType:       input.TriggerType,
		TriggerNegate:     input.TriggerNegate,
		TriggerIgnoreCase: input.TriggerIgnoreCase,
		Description:       input.Description,
		Position:          input.Position,
	})
	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
