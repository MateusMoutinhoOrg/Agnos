package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setBodyAction.SetBody(sandbox, api.RouteBodyProps{
		Path:        props.Path,
		Route:       entries.Route,
		Type:        entries.Type,
		Required:    entries.Required,
		Optional:    entries.Optional,
		MaxBytes:    entries.MaxBytes,
		ContentType: entries.ContentType,
		DropSchema:  entries.DropSchema,
	})
	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
