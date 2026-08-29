package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/utils/parsables"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/utils/smartio"
)

// goVersion is the go directive written into a generated go.mod.
const goVersion = "1.25.0"

// StartFactory returns the closure that fills lib.CoreApi.Start, scaffolding the
// project config directory with one empty yaml per parsable and, when a module
// is given, a go.mod at the project root. Every write goes through a SmartIO
// transaction, so nothing touches the disk unless all of them succeed.
func StartFactory(sandbox *lib.SandBox) func(props lib.StartProps) error {

	return func(props lib.StartProps) error {

		io := smartio.NewSmartIO(sandbox, props.Path)

		write := io.WriteFile
		if props.Force {
			write = io.WriteFileOverwrite
		}

		configDir := props.Path + "/" + sandbox.ProjectName + "Config"
		io.CreateDir(configDir)

		project_conf := parsables.NewProjectConfEmpty(sandbox)
		project_conf.Name = props.ProjectName

		themes_conf := parsables.NewThemesConfEmpty(sandbox)
		themes_conf.AddTheme("LibUsage", "lib-usage", "Documentation explaning how to use the lib")
		themes_conf.AddTheme("Development", "development", "Documentation explaning how to to. build the project, and how to modify the project ")
		ignorable_conf := parsables.NewIgnorableConfEmpty(sandbox)
		path_replacer_conf := parsables.NewPathReplacerConfEmpty(sandbox)

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
			module_conf := parsables.NewModuleConfEmpty(sandbox)
			module_conf.Module = *props.Module
			module_conf.GoVersion = goVersion

			write_error := write(props.Path+"/go.mod", []byte(module_conf.Render()))
			if write_error != nil {
				return write_error
			}
		}

		persist_error := io.Persist()
		if persist_error != nil {
			return persist_error
		}

		sandbox.Deps.Printf("started with path %s \n", props.Path)
		return nil
	}
}
