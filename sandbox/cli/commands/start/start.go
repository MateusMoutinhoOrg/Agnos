package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{

		ValidStartIdentifiers: []string{"start"},
		Category:              "Core Commands",

		Args: []api.CliArg{
			api.CliArg{
				Id:          "path",
				Description: "the dir to start the project",
				Examples: []string{
					sandbox.ProjectName + " start . ",
				},
				Defaults:        []string{"."},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []api.Cliflag{
			api.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.ProjectName + " start -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			api.Cliflag{
				Id:               "force",
				ValidIdentifiers: []string{"--force", "-f"},
				Description:      "Forces the creation of the project, overwriting existing files",
				Examples: []string{
					sandbox.ProjectName + " start -f",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			api.Cliflag{
				Id:               "module",
				ValidIdentifiers: []string{"--module", "-m"},
				Description:      "Module name for go.mod",
				Examples: []string{
					sandbox.ProjectName + " start -m github.com/user/project",
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
			sandbox.ProjectName + " start",
			sandbox.ProjectName + " start .",
			sandbox.ProjectName + " start ./my-project",
			sandbox.ProjectName + " start -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	forceFlag := entries.GetFlagById("force")
	moduleFlag := entries.GetFlagById("module")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	var module *string
	if moduleFlag.Exist && len(moduleFlag.Values) > 0 {
		modVal := moduleFlag.Values[0].String()
		module = &modVal
	}

	start_error := sandbox.Start(api.StartProps{
		Path:        path,
		ProjectName: sandbox.ProjectName,
		Module:      module,
		Force:       forceFlag.Exist,
	})

	if !quietFlag.Exist && start_error != nil {
		sandbox.Deps.Error(start_error.Error())
	}
	if start_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
