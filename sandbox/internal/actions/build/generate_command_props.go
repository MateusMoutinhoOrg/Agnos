package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandPropsDest is where api.CommandProps is declared: beside the
// contracts, since every command package names it and none of them may import
// another.
const commandPropsDest = "sandbox/api/" + utils.CommandPropsFile

// retiredCliFiles are generated files an older build wrote that nothing
// renders any more: the dispatch before it moved to generated/cli/cli, and the
// entries.yaml + handler.go the build wrote for help and version before
// command.yaml. Each one names a symbol the current contract dropped, so a
// tree still carrying it would not compile; it is removed on every build.
var retiredCliFiles = []string{
	utils.GeneratedDir + "/cli/climain.go",
	utils.GeneratedDir + "/cli/new.go",
	"sandbox/internal/commands/help/entries.yaml",
	"sandbox/internal/commands/help/handler.go",
	"sandbox/internal/commands/version/entries.yaml",
	"sandbox/internal/commands/version/handler.go",
}

// GenerateCommandProps renders assets/templates/commandprops.go into
// sandbox/api/commandprops.go — the api.CommandProps one command line's chain
// of commands shares, handed to every InternalPureHandler as its first
// argument.
//
// It is written **once**, like the Handle* files: what a command line carries
// from a middleware to the command after it is the project's to type, so a
// struct already on disk is left as it is.
func GenerateCommandProps(sandbox *api.Sandbox, io *smartio.SmartIO, module string) error {
	for _, retired := range retiredCliFiles {
		if io.IsFile(retired) {
			io.RemoveDir(retired)
		}
	}

	if io.IsFile(commandPropsDest) {
		return nil
	}

	vars := map[string]any{
		"Module":        module,
		"GeneratorName": generatorName(sandbox),
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/commandprops.go", vars, commandPropsDest)
}
