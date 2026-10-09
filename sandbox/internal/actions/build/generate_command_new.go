package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateCommandNew renders assets/templates/command_new.go and
// command_input.go once per command into the directory it sits in under
// sandbox/internal/commands:
// generated.new.go, the api.Command that package declares, derived from its
// command.yaml, which sandbox/internal/cli/generated.new.go collects into
// Cli.Commands; and generated.input.go, the Input its Handle is handed.
func GenerateCommandNew(sandbox *api.Sandbox, io *stagedfs.StagedFS, commands []map[string]any, module string) error {
	for _, command := range commands {
		name, _ := command["CommandName"].(string)
		if name == "" {
			continue
		}
		vars := map[string]any{"Module": module, "GeneratorName": generatorName(sandbox)}
		for key, value := range command {
			vars[key] = value
		}
		dir, _ := command["Dir"].(string)
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/command_new.go", vars, dir+"/"+utils.GeneratedFile(sandbox, utils.UnitNewFile)); err != nil {
			return err
		}
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/command_input.go", vars, dir+"/"+utils.GeneratedFile(sandbox, utils.UnitInputFile)); err != nil {
			return err
		}
	}
	return nil
}
