package add_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_body_field"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	add_error := addBodyFieldAction.AddBodyField(sandbox, api.RouteBodyFieldProps{
		Path:                   entries.Path,
		Route:                  entries.Route,
		Name:                   entries.Name,
		Type:                   entries.Type,
		Required:               entries.Required,
		Array:                  entries.Array,
		Min:                    entries.Min,
		Max:                    entries.Max,
		ExclusiveMin:           entries.ExclusiveMin,
		ExclusiveMax:           entries.ExclusiveMax,
		Format:                 entries.Format,
		Pattern:                entries.Pattern,
		Enum:                   entries.Enum,
		Const:                  entries.Const,
		Nullable:               entries.Nullable,
		MinItems:               entries.MinItems,
		MaxItems:               entries.MaxItems,
		UniqueItems:            entries.UniqueItems,
		AdditionalProperties:   entries.AdditionalProperties,
		NoAdditionalProperties: entries.NoAdditionalProperties,
	})
	if add_error != nil {
		sandbox.Deps.Std.Error("%s\n", add_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
