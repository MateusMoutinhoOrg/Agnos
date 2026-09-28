package add_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_flag"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addFlagAction.AddFlag(sandbox, api.FlagProps{
		Path:        props.Path,
		Command:     entries.Target,
		Name:        entries.Name,
		Keys:        entries.Identifier,
		Description: entries.Description,
		Type:        entries.Type,
		Default:     entries.Default,
		Required:    entries.Required,
		Array:       entries.Array,
		Min:         entries.Min,
		Max:         entries.Max,
		Position:    entries.Position,
	})
	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
