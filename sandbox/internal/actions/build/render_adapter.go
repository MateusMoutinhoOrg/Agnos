package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Render_adapters_availables_standard_new_go(deps *deps.Deps, io *smartio.SmartIO, module string) error {

	vars := map[string]string{
		"Module": module,
	}
	err := utils.RenderTemplateToDest(deps, io, "adapters/availables/standard/new.go", vars, "adapters/availables/standard/new.go")
	if err != nil {
		return err
	}
	return nil
}
