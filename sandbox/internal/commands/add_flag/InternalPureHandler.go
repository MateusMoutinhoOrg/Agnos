package add_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addFlagAction.AddFlag(sandbox, api.FlagProps{
		Path:              props.Path,
		Command:           entries.Target,
		Name:              entries.Name,
		Keys:              entries.Key,
		Type:              entries.Type,
		Required:          entries.Required,
		Default:           entries.Default,
		Min:               entries.Min,
		Max:               entries.Max,
		Enum:              entries.Enum,
		Pattern:           entries.Pattern,
		Trigger:           entries.Trigger,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Description:       entries.Description,
		Position:          entries.Position,
	})
	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
