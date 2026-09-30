package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveCommandInternal deletes every file of the command's directory — in whatever
// folder it sits — plus the directory itself, and every folder the removal
// leaves empty. One holding another command below it is refused, and so is the
// command the build writes itself: it is rendered, not declared. name is the
// command's package name or one of its verbs.
func RemoveCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	name = utils.ResolveCommandName(sandbox, io, name)
	if err := utils.ValidateCommandName(sandbox, name); err != nil {
		return err
	}
	if utils.IsGeneratedCommand(sandbox, name) {
		return sandbox.Deps.Std.Errorf("the %s command is generated and cannot be removed", utils.CommandIdentifier(sandbox, name))
	}

	dir := utils.CommandDir(sandbox, io, name)
	if !io.IsDir(dir) {
		return sandbox.Deps.Std.Errorf("command %q not found", utils.CommandIdentifier(sandbox, name))
	}

	if utils.HoldsOtherUnit(sandbox, io, dir, utils.CommandConfFile) {
		return sandbox.Deps.Std.Errorf("command %q holds another command under %s: move or remove that one first", utils.CommandIdentifier(sandbox, name), dir)
	}

	sandbox.Deps.Std.Log("remove-command removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	utils.PruneEmptyGroups(sandbox, io, utils.CommandsDir, dir)
	return nil
}
