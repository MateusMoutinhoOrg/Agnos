package add_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddFlagInternal parses the target command's entries.yaml, appends the new
// flag (refusing a duplicate name or identifier) and writes the file back.
// When no --identifier is given the flag answers to "--<name>".
func AddFlagInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.FieldProps) error {
	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}

	field, err := utils.NewField(sandbox, props)
	if err != nil {
		return err
	}
	if len(field.Identifiers) == 0 {
		field.Identifiers = []string{"--" + field.Key}
	}
	for _, id := range field.Identifiers {
		if !sandbox.Deps.Stringsdeps.HasPrefix(id, "-") {
			return sandbox.Deps.Std.Errorf("flag identifier %q must start with - or --", id)
		}
	}

	if utils.FindField(sandbox, conf.Flags, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("command %q already has a flag named %q", props.Command, field.Key)
	}
	if utils.FindField(sandbox, conf.Args, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("command %q already has an arg named %q", props.Command, field.Key)
	}
	for _, existing := range conf.Flags {
		for _, id := range existing.Identifiers {
			for _, candidate := range field.Identifiers {
				if id == candidate {
					return sandbox.Deps.Std.Errorf("identifier %q is already used by flag %q", id, existing.Key)
				}
			}
		}
	}

	position, err := utils.CheckPosition(sandbox, "flag", props.Position, conf.Flags)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-flag adding %s to %s \n", field.Key, utils.CommandEntriesPath(sandbox, props.Command))

	conf.Flags = utils.InsertField(conf.Flags, field, position)
	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}
