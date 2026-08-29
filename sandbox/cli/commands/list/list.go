package list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func NewCommand(sandbox *lib.SandBox) lib.CliCommand {
	return lib.CliCommand{

		ValidStartIdentifiers: []string{"list"},
		Category:              "Extensions",

		Args: []lib.CliArg{
			lib.CliArg{
				Id:          "path",
				Description: "the dir of the project to list the extensions of",
				Examples: []string{
					sandbox.Config.ProjectName + " list . ",
				},
				Defaults:        []string{"."},
				RequiredType:    lib.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []lib.Cliflag{
			lib.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.Config.ProjectName + " list -q",
				},
				Type:             lib.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},

		Description:     "List all available extensions",
		LongDescription: "Lists every extension available to the project, showing the ones\nalready installed. If no path is provided, the current directory\nis used.",
		Examples: []string{
			sandbox.Config.ProjectName + " list",
			sandbox.Config.ProjectName + " list .",
			sandbox.Config.ProjectName + " list ./my-project",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *lib.SandBox, entries lib.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")

	list_error := sandbox.Core.List(lib.ListProps{
		Path: pathArg.Values[0].String(),
	})

	if !quietFlag.Exist && list_error != nil {
		sandbox.Deps.Error(list_error.Error())
	}
	if list_error != nil {
		return lib.ExitFailure
	}
	return lib.ExitOk
}
