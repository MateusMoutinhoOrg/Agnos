package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Build(deps *deps.Deps, path string) error {
	deps.Printf("build started with path %s \n", path)

	vars := map[string]string{
		"Name": "Agnos",
	}

	err := utils.RenderTemplateToDest(deps, "sandbox/new.go", vars, "sandbox/new.go")
	if err != nil {
		return err
	}

	deps.Printf("successfully rendered template")
	return nil
}
