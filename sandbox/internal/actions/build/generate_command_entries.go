package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// GenerateCommandEntries renders assets/templates/entries.go once per command
// into sandbox/internal/commands/<name>/entries.go — the typed struct the
// user's handler.go receives, derived from that command's entries.yaml.
func GenerateCommandEntries(deps *deps.Deps, io *smartio.SmartIO, commands []map[string]any) error {
	for _, command := range commands {
		name, _ := command["Name"].(string)
		if name == "" {
			continue
		}
		dest := "sandbox/internal/commands/" + name + "/entries.go"
		if err := utils.RenderTemplateToDest(deps, io, "templates/entries.go", command, dest); err != nil {
			return err
		}
	}
	return nil
}
