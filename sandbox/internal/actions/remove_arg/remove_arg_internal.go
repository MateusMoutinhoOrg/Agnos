package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveArgInternal parses the target command's entries.yaml, drops the
// positional arg called name and writes the file back. The args that
// followed it shift one position up.
func RemoveArgInternal(sandbox *api.Sandbox, io *smartio.SmartIO, command string, name string) error {
	conf, err := utils.LoadCommandConf(sandbox, io, command)
	if err != nil {
		return err
	}

	index := utils.FindField(sandbox, conf.Args, name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("command %q has no arg named %q", command, name)
	}

	sandbox.Deps.Std.Log("remove-arg removing %s from %s \n", conf.Args[index].Key, utils.CommandEntriesPath(sandbox, command))

	conf.Args = utils.RemoveField(conf.Args, index)
	return utils.SaveCommandConf(sandbox, io, command, conf)
}
