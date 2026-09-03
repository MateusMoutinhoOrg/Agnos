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
// The file is only ever *created*: once it exists it is left untouched, so a
// project can edit help's declaration like any other command's.
//
// It must run before CollectCommands, so the declaration is already in the
// transaction when the collector reads it and help flows through the same
// entries.go / climain.go / help-metadata generation as any other command.
func GenerateHelpEntriesYaml(deps *deps.Deps, io *smartio.SmartIO, vars map[string]interface{}) error {
	if io.Exist("sandbox/internal/commands/help/entries.yaml") {
		return nil
	}

	io.CreateDir("sandbox/internal/commands/help")
	return utils.RenderTemplateToDest(
		deps, io,
		"templates/help_entries.yaml",
		vars,
		"sandbox/internal/commands/help/entries.yaml",
	)
}
