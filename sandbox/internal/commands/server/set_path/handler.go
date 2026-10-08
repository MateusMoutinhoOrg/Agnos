package set_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	err := setPathAction.SetPath(sandbox, api.SetPathProps{
		Path:              props.Path,
		Route:             input.Route,
		Name:              input.Name,
		Rename:            input.Rename,
		Start:             input.Start,
		End:               input.End,
		TriggerType:       input.TriggerType,
		TriggerNegate:     input.TriggerNegate,
		TriggerIgnoreCase: input.TriggerIgnoreCase,
		Type:              input.Type,
		Trigger:           input.Trigger,
		Description:       input.Description,
		Clear:             input.Clear,
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
