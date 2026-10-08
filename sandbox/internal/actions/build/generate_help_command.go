package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// helpCommandYaml is the declaration of the help command.
const helpCommandYaml = "sandbox/internal/commands/info/help/command.yaml"

// GenerateHelpCommandYaml renders assets/templates/help_command.yaml into
// sandbox/internal/commands/info/help/command.yaml. The help command is declared
// exactly like every other command; the only difference is that the build
// writes its command.yaml instead of the user.
//
// The file is only ever *created*: once it exists it is left untouched, so a
// project can edit help's declaration like any other command's.
//
// It must run before CollectCommands, so the declaration is already in the
// transaction when the collector reads it and help flows through the same
// new.go / Cli.Commands generation as any other command.
func GenerateHelpCommandYaml(sandbox *api.Sandbox, io *stagedfs.StagedFS, vars map[string]interface{}) error {
	if io.Exists(helpCommandYaml) {
		return nil
	}

	io.CreateDir("sandbox/internal/commands/info/help")
	return utils.RenderTemplateToDest(sandbox, io, "templates/help_command.yaml", vars, helpCommandYaml)
}
