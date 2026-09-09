package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	setBodyAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	set_error := setBodyAction.SetBody(deps, api.RouteBodyProps{
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
		deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
