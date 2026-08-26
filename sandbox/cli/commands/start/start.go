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
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	start_error := sandbox.Start(path, sandbox.ProjectName, "default/module")

	if !quietFlag.Exist && start_error != nil {
		sandbox.Deps.Error(start_error.Error())
	}
	if start_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
