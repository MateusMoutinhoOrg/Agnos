package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// AddArgInternal parses the target command's entries.yaml, inserts the new
// positional arg (at --position, else at the end) and writes the file back.
// Positional args carry no identifiers and bind by their written order, so
// an array arg must stay last.
func AddArgInternal(deps *deps.Deps, io *smartio.SmartIO, props api.FieldProps) error {
	conf, err := utils.LoadCommandConf(deps, io, props.Command)
	if err != nil {
		return err
	}

	props.Identifiers = nil
	field, err := utils.NewField(deps, props)
	if err != nil {
		return err
	}
	if field.Type == "boolean" {
		return deps.Std.Errorf("a positional arg cannot be boolean")
	}

	if utils.FindField(conf.Args, field.Key) >= 0 {
		return deps.Std.Errorf("command %q already has an arg named %q", props.Command, field.Key)
	}
	if utils.FindField(conf.Flags, field.Key) >= 0 {
		return deps.Std.Errorf("command %q already has a flag named %q", props.Command, field.Key)
	}

	position, err := utils.CheckPosition(deps, "arg", props.Position, conf.Args)
	if err != nil {
		return err
	}
	for i, existing := range conf.Args {
		if existing.Array && i < position {
			return deps.Std.Errorf("arg %q is an array and must stay last; insert before it with --position %d", existing.Key, i)
		}
	}
	if field.Array && position != len(conf.Args) {
		return deps.Std.Errorf("an array arg must be the last positional arg")
	}

	deps.Std.Log("add-arg adding %s to %s \n", field.Key, utils.CommandEntriesPath(props.Command))

	conf.Args = utils.InsertField(conf.Args, field, position)
	return utils.SaveCommandConf(deps, io, props.Command, conf)
}
