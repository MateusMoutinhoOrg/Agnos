package set_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_flag"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setFlagAction.SetFlag(sandbox, api.FlagEditProps{
		Path:              props.Path,
		Command:           entries.Target,
		Name:              entries.Name,
		Rename:            entries.Rename,
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
		Clear:             entries.Clear,
	})
	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
