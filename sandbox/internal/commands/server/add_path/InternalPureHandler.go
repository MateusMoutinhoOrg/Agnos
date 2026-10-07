package add_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := addPathAction.AddPath(sandbox, api.RoutePathProps{
		Path:              props.Path,
		Route:             entries.Route,
		Id:                entries.Id,
		Start:             entries.Start,
		End:               entries.End,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Type:              entries.Type,
		Trigger:           entries.Trigger,
		Description:       entries.Description,
		Position:          entries.Position,
	})
	if err != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
