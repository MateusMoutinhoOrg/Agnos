package start

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/rundeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// fallbackGoVersion is the go directive written into a generated go.mod when
// the toolchain cannot be asked which version it is.
const fallbackGoVersion = "1.25.0"

// goVersion asks the installed toolchain for its own version, so a generated
// go.mod does not pin a release older than the compiler that will build it.
// `go env GOVERSION` answers "go1.26.0"; anything unexpected (no toolchain,
// no such directory yet) falls back to the constant above.
func goVersion(deps *deps.Deps, path string) string {
	result, err := deps.Rundeps.Run(rundeps.RunProps{
		Dir:     path,
		Program: "go",
		Args:    []string{"env", "GOVERSION"},
	})
	if err != nil || result.ExitCode != 0 {
		return fallbackGoVersion
	}

	version := strings.TrimPrefix(strings.TrimSpace(result.Output), "go")
	if version == "" {
		return fallbackGoVersion
	}
	return version
}

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
		module_conf.GoVersion = goVersion(deps, props.Path)

		if err := write("go.mod", []byte(module_conf.Render())); err != nil {
			return err
		}
	}

	deps.Std.Log("started with path %s \n", props.Path)
	return nil
}
