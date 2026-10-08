package set_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_flag"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setFlagAction.SetFlag(sandbox, api.SetFlagProps{
		Path:              props.Path,
		Command:           input.Command,
		Name:              input.Name,
		Rename:            input.Rename,
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
		Clear:             input.Clear,
	})
	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
