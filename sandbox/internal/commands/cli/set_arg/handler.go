package set_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setArgAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_arg"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setArgAction.SetArg(sandbox, api.SetArgProps{
		Path:              props.Path,
		Command:           input.Command,
		Name:              input.Name,
		Rename:            input.Rename,
		Start:             input.Start,
		End:               input.End,
		Type:              input.Type,
		Required:          input.Required,
		Default:           input.Default,
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
