package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Build(deps *deps.Deps, path string) error {
	deps.Printf("build started with path %s \n", path)

	io := smartio.New(deps, path, config.ProjectName)

	vars := map[string]string{
		"Name": "Agnos",
	}

	err := utils.RenderTemplateToDest(deps, io, "sandbox/new.go", vars, "sandbox/new.go")
	if err != nil {
		return err
	}

	err = io.Persist()
	if err != nil {
		return err
	}

	deps.Printf("successfully rendered template")
	return nil
}
