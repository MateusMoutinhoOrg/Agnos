package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setCommandAction.SetCommand(sandbox, api.SetCommandProps{
		Path:        props.Path,
		Command:     input.Name,
		Summary:     input.Summary,
		Category:    input.Category,
		Description: input.Description,
		Hidden:      input.Hidden,
		Visible:     input.Visible,
		Strict:      input.Strict,
		Loose:       input.Loose,
		Priority:    input.Priority,
		HasPriority: input.Priority >= 0,
		Before:      input.Before,
		After:       input.After,
		Segments:    input.Segments,
		HasSegments: input.Segments >= 0,
		Identifiers: input.Identifier,
		Examples:    input.Example,
		Clear:       input.Clear,
	})
	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
