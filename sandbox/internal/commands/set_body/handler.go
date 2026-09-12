package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	set_error := setBodyAction.SetBody(sandbox, api.RouteBodyProps{
		Path:        entries.Path,
		Route:       entries.Route,
		Type:        entries.Type,
		Required:    entries.Required,
		Optional:    entries.Optional,
		MaxBytes:    entries.MaxBytes,
		ContentType: entries.ContentType,
		DropSchema:  entries.DropSchema,
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
