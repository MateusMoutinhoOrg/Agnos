package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// RemoveCommandInternal deletes every file under
// sandbox/internal/commands/<name>/ plus the directory itself. The generated
// help command is refused: it is rendered by build, not declared.
func RemoveCommandInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateCommandName(deps, name); err != nil {
		return err
	}
	pkg := utils.CommandPackage(name)
	if pkg == "help" {
		return deps.Std.Errorf("the help command is generated and cannot be removed")
	}

	dir := utils.CommandDir(name)
	if !io.IsDir(dir) {
		return deps.Std.Errorf("command %q not found", utils.CommandIdentifier(name))
	}

	deps.Std.Log("remove-command removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
