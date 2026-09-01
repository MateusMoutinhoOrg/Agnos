package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func Build(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Printf("build started with path %s \n", path)

	//Creating the basic dir struct
	io.CreateDir("sandbox/api")
	io.CreateDir("sandbox/internal")
	/*
		project_conf, err := projectconf.New(deps, path)
		if err != nil {
			return err
		}
	*/
	module_conf, err := moduleconf.NewFromPath(deps, path+"/go.mod")
	if err != nil {
		return err
	}
	err = Render_sandbox_new_go(deps, io, module_conf.Module)
	if err != nil {
		return err
	}

	err = Render_sandbox_api_sandbox_go(deps, io)
	if err != nil {
		return err
	}

	deps.Printf("successfully rendered template")
	return nil
}
