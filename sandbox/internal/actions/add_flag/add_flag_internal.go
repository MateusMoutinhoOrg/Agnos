package add_flag

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddFlagInternal parses the target command's entries.yaml, appends the new
// flag (refusing a duplicate name or identifier) and writes the file back.
// When no --identifier is given the flag answers to "--<name>".
func AddFlagInternal(deps *deps.Deps, io *smartio.SmartIO, props api.FieldProps) error {
	conf, err := utils.LoadCommandConf(deps, io, props.Command)
	if err != nil {
		return err
	}

	field, err := utils.NewField(deps, props)
	if err != nil {
		return err
	}
	if len(field.Identifiers) == 0 {
		field.Identifiers = []string{"--" + field.Key}
	}
	for _, id := range field.Identifiers {
		if !strings.HasPrefix(id, "-") {
			return deps.Std.Errorf("flag identifier %q must start with - or --", id)
		}
	}

	if utils.FindField(conf.Flags, field.Key) >= 0 {
		return deps.Std.Errorf("command %q already has a flag named %q", props.Command, field.Key)
	}
	if utils.FindField(conf.Args, field.Key) >= 0 {
		return deps.Std.Errorf("command %q already has an arg named %q", props.Command, field.Key)
	}
	for _, existing := range conf.Flags {
		for _, id := range existing.Identifiers {
			for _, candidate := range field.Identifiers {
				if id == candidate {
					return deps.Std.Errorf("identifier %q is already used by flag %q", id, existing.Key)
				}
			}
		}
	}

	position, err := utils.CheckPosition(deps, "flag", props.Position, conf.Flags)
	if err != nil {
		return err
	}

	deps.Std.Log("add-flag adding %s to %s \n", field.Key, utils.CommandEntriesPath(props.Command))

	conf.Flags = utils.InsertField(conf.Flags, field, position)
	return utils.SaveCommandConf(deps, io, props.Command, conf)
}
