package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addCommandAction.AddCommand(sandbox, api.AddCommandProps{
		Path:              props.Path,
		Name:              entries.Name,
		Trigger:           entries.Trigger,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Pattern:           entries.Pattern,
		Middleware:        entries.Middleware,
		Priority:          entries.Priority,
		HasPriority:       entries.Priority >= 0,
		Before:            entries.Before,
		After:             entries.After,
		Help:              entries.Help,
		Category:          entries.Category,
	})
	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
