package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveArgInternal parses the target command's command.yaml, drops the arg
// called name and writes the file back. The segments it read are left for no
// arg: the args after it keep theirs.
func RemoveArgInternal(sandbox *api.Sandbox, io *smartio.SmartIO, command string, name string) error {
	conf, err := utils.LoadCommandConf(sandbox, io, command)
	if err != nil {
		return err
	}

	index := utils.FindCommandArg(conf, utils.CommandEntryId(sandbox, name))
	if index < 0 {
		return sandbox.Deps.Std.Errorf("command %q has no arg named %q", command, name)
	}
	if len(conf.Args) == 1 {
		return sandbox.Deps.Std.Errorf("arg %q is the only one of command %q: a command reads one segment at least", name, command)
	}

	sandbox.Deps.Std.Log("remove-arg removing %s from %s \n", conf.Args[index].Id, utils.CommandConfPath(sandbox, utils.ResolveCommandName(sandbox, io, command)))

	conf.Args = utils.RemoveCommandArg(conf.Args, index)
	return utils.SaveCommandConf(sandbox, io, command, conf)
}
