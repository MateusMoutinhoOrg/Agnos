package add_segment

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	addSegmentAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_segment"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	add_error := addSegmentAction.AddSegment(deps, api.RouteFieldProps{
		Path:        entries.Path,
		Route:       entries.Route,
		Name:        entries.Name,
		Identifier:  entries.Identifier,
		Description: entries.Description,
		Examples:    entries.Example,
		Type:        entries.Type,
		Array:       entries.Array,
		Min:         entries.Min,
		Max:         entries.Max,
		Position:    entries.Position,
	})
	if add_error != nil {
		deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
