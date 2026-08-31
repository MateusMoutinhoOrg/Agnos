package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"start"},
		Category:              "Core Commands",
		Args: []api.CliArg{
			{
				Id:          "project-name",
				Description: "the name of the project",
				Examples: []string{
					"my-project",
				},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
			{
				Id:          "path",
				Description: "the dir to start the project",
				Examples: []string{
					".",
				},
				Defaults:        []string{"."},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []api.Cliflag{
			{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					config.ProjectName + " start -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			{
				Id:               "force",
				ValidIdentifiers: []string{"--force", "-f"},
				Description:      "Forces the creation of the project, overwriting existing files",
				Examples: []string{
					config.ProjectName + " start -f",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			{
				Id:               "module",
				ValidIdentifiers: []string{"--module", "-m"},
				Description:      "Module name for go.mod",
				Examples: []string{
					config.ProjectName + " start -m github.com/user/project",
				},
				Type:             api.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},
		Description:     "Initialize a new project in a directory",
		LongDescription: "Scaffolds a new Agnos project in the given directory, creating\nthe required configuration files and folder structure. If no\npath is provided, the current directory is used.",
		Examples: []string{
			config.ProjectName + " start my-project",
			config.ProjectName + " start my-project .",
			config.ProjectName + " start my-project ./my-project-dir",
			config.ProjectName + " start my-project -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	forceFlag := entries.GetFlagById("force")
	moduleFlag := entries.GetFlagById("module")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()
	projectNameArg := entries.GetArgById("project-name")
	projectName := projectNameArg.Values[0].String()

	var module *string
	if moduleFlag.Exist && len(moduleFlag.Values) > 0 {
		modVal := moduleFlag.Values[0].String()
		module = &modVal
	}

	if !deps.IoLib.Exist(path + "/go.mod") && module == nil {
		if !quietFlag.Exist {
			deps.Error("the module flag (--module) is required when there is no go.mod in the path\n")
		}
		return api.ExitFailure
	}

	start_error := startAction.Start(deps, api.StartProps{
		Path:        path,
		ProjectName: projectName,
		Module:      module,
		Force:       forceFlag.Exist,
	})

	if !quietFlag.Exist && start_error != nil {
		deps.Error(start_error.Error())
	}
	if start_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
