package set_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_body_field"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	set_error := setBodyFieldAction.SetBodyField(sandbox, api.RouteBodyFieldEditProps{
		Path:                   command.GetString("path"),
		Route:                  command.GetString("route"),
		Name:                   command.GetString("name"),
		Rename:                 command.GetString("rename"),
		Type:                   command.GetString("type"),
		Required:               command.GetBool("required"),
		Array:                  command.GetBool("array"),
		Min:                    command.GetString("min"),
		Max:                    command.GetString("max"),
		ExclusiveMin:           command.GetString("exclusive-min"),
		ExclusiveMax:           command.GetString("exclusive-max"),
		Format:                 command.GetString("format"),
		Pattern:                command.GetString("pattern"),
		Enum:                   command.GetStrings("enum"),
		Const:                  command.GetString("const"),
		Nullable:               command.GetBool("nullable"),
		MinItems:               command.GetString("min-items"),
		MaxItems:               command.GetString("max-items"),
		UniqueItems:            command.GetBool("unique-items"),
		AdditionalProperties:   command.GetBool("additional-properties"),
		NoAdditionalProperties: command.GetBool("no-additional-properties"),
		Clear:                  command.GetStrings("clear"),
	})
	if set_error != nil {
		sandbox.Deps.Std.Error("%s\n", set_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
