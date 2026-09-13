package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateCommandNew renders assets/templates/new.go once per command into
// sandbox/internal/commands/<name>/new.go — the api.Command that package
// declares, derived from its entries.yaml, which sandbox/binds/cli.go collects
// into sandbox.Commands.
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
		dest := "sandbox/internal/commands/" + name + "/new.go"
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/new.go", vars, dest); err != nil {
			return err
		}
	}
	return nil
}
