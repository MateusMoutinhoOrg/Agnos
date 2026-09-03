package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// RemoveArgInternal parses the target command's entries.yaml, drops the
// positional arg called name and writes the file back. The args that
// followed it shift one position up.
func RemoveArgInternal(deps *deps.Deps, io *smartio.SmartIO, command string, name string) error {
	conf, err := utils.LoadCommandConf(deps, io, command)
	if err != nil {
		return err
	}

	index := utils.FindField(conf.Args, name)
	if index < 0 {
		return deps.Std.Errorf("command %q has no arg named %q", command, name)
	}

	deps.Std.Log("remove-arg removing %s from %s \n", conf.Args[index].Key, utils.CommandEntriesPath(command))

	conf.Args = utils.RemoveField(conf.Args, index)
	return utils.SaveCommandConf(deps, io, command, conf)
}
