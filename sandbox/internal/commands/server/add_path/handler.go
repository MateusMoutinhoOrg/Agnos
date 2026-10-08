package add_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	err := addPathAction.AddPath(sandbox, api.AddPathProps{
		Path:              props.Path,
		Route:             input.Route,
		Name:              input.Name,
		Start:             input.Start,
		End:               input.End,
		TriggerType:       input.TriggerType,
		TriggerNegate:     input.TriggerNegate,
		TriggerIgnoreCase: input.TriggerIgnoreCase,
		Type:              input.Type,
		Trigger:           input.Trigger,
		Description:       input.Description,
		Position:          input.Position,
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
