package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/ignorableconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/pathreplacerconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/themesconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// goVersion is the go directive written into a generated go.mod.
const goVersion = "1.25.0"

func Start(deps *deps.Deps, props api.StartProps) error {

	io := smartio.New(deps, props.Path, props.ProjectName)

	write := io.WriteFile
	if props.Force {
		write = io.WriteFileOverwrite
	}

	configDir := props.Path + "/" + config.ProjectName + "Config"
	io.CreateDir(configDir)

	project_conf := projectconf.NewEmpty(deps)
	project_conf.Name = props.ProjectName

	themes_conf := themesconf.NewEmpty(deps)
	themes_conf.AddTheme("LibUsage", "lib-usage", "Documentation explaning how to use the lib")
	themes_conf.AddTheme("Development", "development", "Documentation explaning how to to. build the project, and how to modify the project ")
	ignorable_conf := ignorableconf.NewEmpty(deps)
	path_replacer_conf := pathreplacerconf.NewEmpty(deps)

	parsable_files := []struct {
		path    string
		content string
	}{
		{configDir + "/project.yaml", project_conf.Render()},
		{configDir + "/themes.yaml", themes_conf.Render()},
		{configDir + "/ignore.yaml", ignorable_conf.Render()},
		{configDir + "/paths.yaml", path_replacer_conf.Render()},
	}

	for _, file := range parsable_files {
		write_error := write(file.path, []byte(file.content))
		if write_error != nil {
			return write_error
		}
	}

	if props.Module != nil {
		module_conf := moduleconf.NewEmpty(deps)
		module_conf.Module = *props.Module
		module_conf.GoVersion = goVersion

		write_error := write(props.Path+"/go.mod", []byte(module_conf.Render()))
		if write_error != nil {
			return write_error
		}
	}

	io.CreateDir("sandbox/api")
	io.CreateDir("sandbox/binds")

	persist_error := io.Persist()
	if persist_error != nil {
		return persist_error
	}
	build.Build(deps, props.Path)
	deps.Printf("started with path %s \n", props.Path)
	return nil
}
