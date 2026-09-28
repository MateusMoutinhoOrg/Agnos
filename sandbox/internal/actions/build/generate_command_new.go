package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateCommandNew renders assets/templates/command_new.go and
// command_entries.go once per command into sandbox/internal/commands/<name>/:
// new.go, the api.Command that package declares, derived from its
// command.yaml, which sandbox/internal/generated/cli/cli/new.go collects into
// Cli.Commands; and entries.go, the Entries its InternalPureHandler is handed.
func GenerateCommandNew(sandbox *api.Sandbox, io *smartio.SmartIO, commands []map[string]any, module string) error {
	for _, command := range commands {
		name, _ := command["Name"].(string)
		if name == "" {
			continue
		}
		vars := map[string]any{"Module": module, "GeneratorName": generatorName(sandbox)}
		for key, value := range command {
			vars[key] = value
		}
		dir := "sandbox/internal/commands/" + name
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/command_new.go", vars, dir+"/new.go"); err != nil {
			return err
		}
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/command_entries.go", vars, dir+"/entries.go"); err != nil {
			return err
		}
	}
	return nil
}
