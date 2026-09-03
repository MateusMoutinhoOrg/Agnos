package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveFlagInternal parses the target command's entries.yaml, drops the
// flag called name (matched by its declared name, or by one of its
// identifiers such as "--out") and writes the file back.
func RemoveFlagInternal(deps *deps.Deps, io *smartio.SmartIO, command string, name string) error {
	conf, err := utils.LoadCommandConf(deps, io, command)
	if err != nil {
		return err
	}

	index := utils.FindField(conf.Flags, name)
	if index < 0 {
		index = findByIdentifier(conf, name)
	}
	if index < 0 {
		return deps.Std.Errorf("command %q has no flag named %q", command, name)
	}

	deps.Std.Log("remove-flag removing %s from %s \n", conf.Flags[index].Key, utils.CommandEntriesPath(command))

	conf.Flags = utils.RemoveField(conf.Flags, index)
	return utils.SaveCommandConf(deps, io, command, conf)
}

// findByIdentifier locates a flag by one of its identifiers ("--out", "-o").
func findByIdentifier(conf *commandconf.CommandConf, name string) int {
	for i, flag := range conf.Flags {
		for _, id := range flag.Identifiers {
			if id == name {
				return i
			}
		}
	}
	return -1
}
