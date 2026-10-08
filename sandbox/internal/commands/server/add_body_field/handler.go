package add_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_body_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addBodyFieldAction.AddBodyField(sandbox, api.AddBodyFieldProps{
		Path:                   props.Path,
		Route:                  input.Route,
		Name:                   input.Name,
		Type:                   input.Type,
		Required:               input.Required,
		Array:                  input.Array,
		Min:                    input.Min,
		Max:                    input.Max,
		ExclusiveMin:           input.ExclusiveMin,
		ExclusiveMax:           input.ExclusiveMax,
		Format:                 input.Format,
		Pattern:                input.Pattern,
		Enum:                   input.Enum,
		Const:                  input.Const,
		Nullable:               input.Nullable,
		MinItems:               input.MinItems,
		MaxItems:               input.MaxItems,
		UniqueItems:            input.UniqueItems,
		AdditionalProperties:   input.AdditionalProperties,
		NoAdditionalProperties: input.NoAdditionalProperties,
	})
	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
