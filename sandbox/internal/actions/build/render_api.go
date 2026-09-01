package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Render_sandbox_api_sandbox_go(deps *deps.Deps, io *smartio.SmartIO) error {
	vars := map[string]string{}
	err := utils.RenderTemplateToDest(deps, io, "sandbox/api/sandbox.go", vars, "sandbox/api/sandbox.go")
	if err != nil {
		return err
	}
	return nil
}
