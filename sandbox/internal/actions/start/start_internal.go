package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// StartInternal scaffolds a project at props.Path. A directory that already
// holds one is refused unless --force is given, and then only the name in
// project.yaml changes: extensions.yaml, ReadmeHeader.md and every other file
// start writes once are the project's by now, and re-rendering them would
// reset what it declared and wrote. A module other than the one go.mod
// already declares is refused, since every import of the tree names the old
// one.
func StartInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.StartProps) error {
	if err := utils.ValidateProjectName(sandbox, props.ProjectName); err != nil {
		return err
	}

	if io.IsFile(utils.ProjectConfPath(sandbox)) {
		return restartInternal(sandbox, io, props)
	}

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

// restartInternal is start over a directory that is already a project: with
// --force it renames it and keeps everything else.
func restartInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.StartProps) error {
	if !props.Force {
		return sandbox.Deps.Std.Errorf("%s already exists: this directory is already a project (pass --force to rename it, keeping every other file)", utils.ProjectConfPath(sandbox))
	}

	if props.Module != nil && io.IsFile("go.mod") {
		module_conf, err := utils.LoadModuleConf(sandbox, io)
		if err != nil {
			return err
		}
		if module_conf.Module != *props.Module {
			return sandbox.Deps.Std.Errorf("the project's module is %s: changing it to %s is not supported, every import of the tree names the current one", module_conf.Module, *props.Module)
		}
	}

	project_conf, err := utils.LoadProjectConf(sandbox, io)
	if err != nil {
		return err
	}
	if project_conf.Name != props.ProjectName {
		sandbox.Deps.Std.Log("start: renaming the project from %q to %q, keeping every other file \n", project_conf.Name, props.ProjectName)
	}
	project_conf.Name = props.ProjectName
	return io.WriteFileOverwrite(utils.ProjectConfPath(sandbox), []byte(project_conf.Render()))
}
