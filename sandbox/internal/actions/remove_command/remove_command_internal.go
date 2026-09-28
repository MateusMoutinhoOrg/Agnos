package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveCommandInternal deletes every file under
// sandbox/internal/commands/<name>/ plus the directory itself, name being the
// command's package name or one of its verbs. The commands the build writes
// itself are refused: they are rendered, not declared.
func RemoveCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	name = utils.ResolveCommandName(sandbox, io, name)
	if err := utils.ValidateCommandName(sandbox, name); err != nil {
		return err
	}
	if utils.IsGeneratedCommand(sandbox, name) {
		return sandbox.Deps.Std.Errorf("the %s command is generated and cannot be removed", utils.CommandIdentifier(sandbox, name))
	}

	dir := utils.CommandDir(sandbox, name)
	if !io.IsDir(dir) {
		return sandbox.Deps.Std.Errorf("command %q not found", utils.CommandIdentifier(sandbox, name))
	}

	sandbox.Deps.Std.Log("remove-command removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
