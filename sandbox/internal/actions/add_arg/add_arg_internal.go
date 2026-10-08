package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddArgInternal parses the target command's command.yaml, inserts the new
// arg (at --position, else at the end) and writes the file back. An arg given
// no --start reads the first segment no arg reads yet.
func AddArgInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddArgProps) error {
	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}

	arg, err := utils.NewCommandArg(sandbox, props, utils.NextCommandSegment(conf))
	if err != nil {
		return err
	}
	if utils.CommandIdTaken(conf, arg.Id) {
		return sandbox.Deps.StdDeps.Errorf("command %q already has an arg or a flag named %s", props.Command, arg.Id)
	}

	position, err := utils.CheckCommandPosition(sandbox, "arg", props.Position, len(conf.Args))
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("add-arg adding %s to %s \n", arg.Id, utils.CommandConfPath(sandbox, io, props.Command))

	conf.Args = utils.InsertAt(conf.Args, arg, position)
	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}
