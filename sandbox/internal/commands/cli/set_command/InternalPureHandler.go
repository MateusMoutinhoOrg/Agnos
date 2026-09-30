package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setCommandAction.SetCommand(sandbox, api.SetCommandProps{
		Path:            props.Path,
		Command:         entries.Name,
		Help:            entries.Help,
		Category:        entries.Category,
		LongDescription: entries.LongDescription,
		Hidden:          entries.Hidden,
		Visible:         entries.Visible,
		Strict:          entries.Strict,
		Loose:           entries.Loose,
		Priority:        entries.Priority,
		HasPriority:     entries.Priority >= 0,
		Before:          entries.Before,
		After:           entries.After,
		Segments:        entries.Segments,
		HasSegments:     entries.Segments >= 0,
		Identifiers:     entries.Identifier,
		Examples:        entries.Example,
		Clear:           entries.Clear,
	})
	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
