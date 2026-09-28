package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveFlagInternal parses the target command's command.yaml, drops the flag
// called name (matched by its name, its id, or one of its keys such as
// "--out") and writes the file back.
func RemoveFlagInternal(sandbox *api.Sandbox, io *smartio.SmartIO, command string, name string) error {
	conf, err := utils.LoadCommandConf(sandbox, io, command)
	if err != nil {
		return err
	}

	index := utils.FindCommandFlagNamed(sandbox, conf, name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("command %q has no flag named %q", command, name)
	}

	sandbox.Deps.Std.Log("remove-flag removing %s from %s \n", conf.Flags[index].Id, utils.CommandConfPath(sandbox, utils.ResolveCommandName(sandbox, io, command)))

	conf.Flags = utils.RemoveCommandFlag(conf.Flags, index)
	return utils.SaveCommandConf(sandbox, io, command, conf)
}
