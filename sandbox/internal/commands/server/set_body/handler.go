package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setBodyAction.SetBody(sandbox, api.SetBodyProps{
		Path:        props.Path,
		Route:       input.Name,
		Type:        input.Type,
		Required:    input.Required,
		Optional:    input.Optional,
		MaxBytes:    input.MaxBytes,
		ContentType: input.ContentType,
		DropSchema:  input.DropSchema,
	})
	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
