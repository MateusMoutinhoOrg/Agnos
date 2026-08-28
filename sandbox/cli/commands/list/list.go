package list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{

		ValidStartIdentifiers: []string{"list"},
		Category:              "Extensions",

		Args: []api.CliArg{
			api.CliArg{
				Id:          "path",
				Description: "the dir of the project to list the extensions of",
				Examples: []string{
					sandbox.ProjectName + " list . ",
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
					sandbox.ProjectName + " list -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},

		Description:     "List all available extensions",
		LongDescription: "Lists every extension available to the project, showing the ones\nalready installed. If no path is provided, the current directory\nis used.",
		Examples: []string{
			sandbox.ProjectName + " list",
			sandbox.ProjectName + " list .",
			sandbox.ProjectName + " list ./my-project",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")

	list_error := sandbox.List(api.ListProps{
		Path: pathArg.Values[0].String(),
	})

	if !quietFlag.Exist && list_error != nil {
		sandbox.Deps.Error(list_error.Error())
	}
	if list_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
