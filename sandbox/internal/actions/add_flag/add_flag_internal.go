package add_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddFlagInternal parses the target command's command.yaml, inserts the new
// flag (refusing an id or a key the command already uses) and writes the file
// back. When no --key is given the flag answers to "--<name>".
func AddFlagInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddFlagProps) error {
	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}

	flag, err := utils.NewCommandFlag(sandbox, props)
	if err != nil {
		return err
	}
	if utils.CommandIdTaken(conf, flag.Id) {
		return sandbox.Deps.StdDeps.Errorf("command %q already has an arg or a flag named %s", props.Command, flag.Id)
	}
	if taken := utils.CommandKeyTaken(conf, flag.Keys); taken != "" {
		return sandbox.Deps.StdDeps.Errorf("key %q is already used by another flag of %q", taken, props.Command)
	}

	position, err := utils.CheckCommandPosition(sandbox, "flag", props.Position, len(conf.Flags))
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("add-flag adding %s to %s \n", flag.Id, utils.CommandConfPath(sandbox, io, props.Command))

	conf.Flags = utils.InsertAt(conf.Flags, flag, position)
	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}
