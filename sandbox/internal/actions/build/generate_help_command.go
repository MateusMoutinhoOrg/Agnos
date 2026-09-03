package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// GenerateHelpEntriesYaml renders assets/templates/help_entries.yaml into
// sandbox/internal/commands/help/entries.yaml. The help command is declared
// exactly like every other command; the only difference is that agnos writes
// its entries.yaml instead of the user.
//
// It must run before CollectCommands, so the declaration is already in the
// transaction when the collector reads it and help flows through the same
// entries.go / climain.go / help-metadata generation as any other command.
func GenerateHelpEntriesYaml(deps *deps.Deps, io *smartio.SmartIO, vars map[string]interface{}) error {
	io.CreateDir("sandbox/internal/commands/help")
	return utils.RenderTemplateToDest(
		deps, io,
		"templates/help_entries.yaml",
		vars,
		"sandbox/internal/commands/help/entries.yaml",
	)
}
