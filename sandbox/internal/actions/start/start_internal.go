package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// goVersion is the go directive written into a generated go.mod.
const goVersion = "1.25.0"

func StartInternal(deps *deps.Deps, io *smartio.SmartIO, props api.StartProps) error {

	project_conf := projectconf.NewEmpty(deps)
	project_conf.Name = props.ProjectName

	vars := map[string]interface{}{
		"Name":        project_conf.Name,
		"Version":     project_conf.Version,
		"Description": project_conf.Description,
	}

	if err := utils.RenderGroup(deps, io, "start", vars); err != nil {
		return err
	}

	if props.Module != nil {
		write := io.WriteFile
		if props.Force {
			write = io.WriteFileOverwrite
		}

		module_conf := moduleconf.NewEmpty(deps)
		module_conf.Module = *props.Module
		module_conf.GoVersion = goVersion

		if err := write(props.Path+"/go.mod", []byte(module_conf.Render())); err != nil {
			return err
		}
	}

	deps.Std.Printf("started with path %s \n", props.Path)
	return nil
}
