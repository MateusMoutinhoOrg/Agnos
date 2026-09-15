package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddArgInternal parses the target command's entries.yaml, inserts the new
// positional arg (at --position, else at the end) and writes the file back.
// Positional args carry no identifiers and bind by their written order, so
// an array arg must stay last.
func AddArgInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.FieldProps) error {
	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}

	props.Identifiers = nil
	field, err := utils.NewField(sandbox, props)
	if err != nil {
		return err
	}
	if field.Type == "boolean" {
		return sandbox.Deps.Std.Errorf("a positional arg cannot be boolean")
	}

	if utils.FindField(sandbox, conf.Args, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("command %q already has an arg named %q", props.Command, field.Key)
	}
	if utils.FindField(sandbox, conf.Flags, field.Key) >= 0 {
		return sandbox.Deps.Std.Errorf("command %q already has a flag named %q", props.Command, field.Key)
	}

	position, err := utils.CheckPosition(sandbox, "arg", props.Position, conf.Args)
	if err != nil {
		return err
	}
	for i, existing := range conf.Args {
		if existing.Array && i < position {
			return sandbox.Deps.Std.Errorf("arg %q is an array and must stay last; insert before it with --position %d", existing.Key, i)
		}
	}
	if field.Array && position != len(conf.Args) {
		return sandbox.Deps.Std.Errorf("an array arg must be the last positional arg")
	}

	sandbox.Deps.Std.Log("add-arg adding %s to %s \n", field.Key, utils.CommandEntriesPath(sandbox, props.Command))

	conf.Args = utils.InsertField(conf.Args, field, position)
	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}
