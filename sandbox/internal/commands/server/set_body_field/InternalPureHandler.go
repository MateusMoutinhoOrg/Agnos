package set_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setBodyFieldAction.SetBodyField(sandbox, api.RouteBodyFieldEditProps{
		Path:                   props.Path,
		Route:                  entries.Route,
		Name:                   entries.Name,
		Rename:                 entries.Rename,
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
		Clear:                  entries.Clear,
	})
	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
