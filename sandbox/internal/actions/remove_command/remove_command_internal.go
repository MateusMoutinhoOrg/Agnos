package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveCommandInternal deletes every file under
// sandbox/internal/commands/<name>/ plus the directory itself. The generated
// help command is refused: it is rendered by build, not declared.
func RemoveCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateCommandName(sandbox, name); err != nil {
		return err
	}
	pkg := utils.CommandPackage(sandbox, name)
	if pkg == "help" {
		return sandbox.Deps.Std.Errorf("the help command is generated and cannot be removed")
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
