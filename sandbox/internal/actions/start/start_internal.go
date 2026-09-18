package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func StartInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.StartProps) error {

	project_conf := projectconf.NewEmpty(sandbox)
	project_conf.Name = props.ProjectName

	vars := map[string]interface{}{
		"Name":      project_conf.Name,
		"Version":   project_conf.Version,
		"ConfigDir": sandbox.Config.ProjectName + "Config",
		"GoRelease": utils.GoRelease,
		"GoFloor":   utils.GoFloor,
	}

	if err := utils.RenderGroup(sandbox, io, "start", vars); err != nil {
		return err
	}

	if props.Module != nil {
		write := io.WriteFile
		if props.Force {
			write = io.WriteFileOverwrite
		}

		module_conf := moduleconf.NewEmpty(sandbox)
		module_conf.Module = *props.Module
		module_conf.GoVersion = utils.GoRelease

		if err := write("go.mod", []byte(module_conf.Render())); err != nil {
			return err
		}
	}

	sandbox.Deps.Std.Log("started with path %s \n", props.Path)
	return nil
}
