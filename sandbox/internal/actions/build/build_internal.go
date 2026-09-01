package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func BuildInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Printf("build started with path %s \n", path)

	//Creating the basic dir struct
	io.CreateDir("sandbox/api")
	io.CreateDir("sandbox/internal")
	/*
		project_conf, err := projectconf.New(deps, path)
		if err != nil {
			return err
		}
	*/
	gomod, err := io.ReadFile(path + "/go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
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
	if io.IsDir("sandbox/deps") {
		err = Render_sandbox_deps_deps_go(deps, io, module_conf.Module)
		if err != nil {
			return err
		}

		err = Render_adapters_availables_standard_new_go(deps, io, module_conf.Module)
		if err != nil {
			return err
		}

	}

	deps.Std.Printf("successfully rendered template")
	return nil
}
