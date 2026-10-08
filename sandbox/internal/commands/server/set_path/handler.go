package set_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setPathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := setPathAction.SetPath(sandbox, api.RoutePathEditProps{
		Path:              props.Path,
		Route:             entries.Route,
		Id:                entries.Id,
		Rename:            entries.Rename,
		Start:             entries.Start,
		End:               entries.End,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Type:              entries.Type,
		Trigger:           entries.Trigger,
		Description:       entries.Description,
		Clear:             entries.Clear,
	})
	if err != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
