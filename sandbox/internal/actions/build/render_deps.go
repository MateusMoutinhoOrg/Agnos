package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Render_sandbox_deps_deps_go(deps *deps.Deps, io *smartio.SmartIO, module string) error {

	vars := map[string]string{
		"Module": module,
	}
	err := utils.RenderTemplateToDest(deps, io, "sandbox/deps/deps.go", vars, "sandbox/deps/deps.go")
	if err != nil {
		return err
	}
	return nil
}
