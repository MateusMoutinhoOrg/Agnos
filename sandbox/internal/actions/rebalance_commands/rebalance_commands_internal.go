package rebalance_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RebalanceCommandsInternal gives every command a rung of its own, Step apart
// and in the order the chain runs them now: the first on Step, the next on
// twice Step, and so on. Two commands that shared a rung ran by name, and
// still do, each on a rung of its own; what changes is that --before and
// --after have room again. The commands the build writes itself keep their
// rungs: the build writes them back.
func RebalanceCommandsInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.RebalanceCommandsProps) error {
	if props.Step < 1 {
		return sandbox.Deps.StdDeps.Errorf("--step %d is below 1: the rungs have to be apart", props.Step)
	}

	chain, err := utils.LoadCommandChain(sandbox, io)
	if err != nil {
		return err
	}

	rung := 0
	for _, entry := range chain {
		if utils.IsGeneratedCommand(sandbox, entry.Name) {
			continue
		}
		rung += props.Step
		if entry.Conf.Priority == rung {
			continue
		}
		sandbox.Deps.StdDeps.Logf("rebalance-commands %s: %d -> %d \n", utils.CommandName(sandbox, entry.Name), entry.Conf.Priority, rung)
		entry.Conf.Priority = rung
		if err := utils.SaveCommandConf(sandbox, io, entry.Name, entry.Conf); err != nil {
			return err
		}
	}
	return nil
}
