package list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{

		ValidStartIdentifiers: []string{"list"},
		Category:              "Extensions",

		Args: []cli.CliArg{
			cli.CliArg{
				Id:          "path",
				Description: "the dir of the project to list the extensions of",
				Examples: []string{
					sandbox.Config.ProjectName + " list . ",
				},
				Defaults:        []string{"."},
				RequiredType:    cli.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []cli.Cliflag{
			cli.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.Config.ProjectName + " list -q",
				},
				Type:             cli.CliTypeBool,
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

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)

	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")

	list_error := sandbox.Core.List(core.ListProps{
		Path: pathArg.Values[0].String(),
	})

	if !quietFlag.Exist && list_error != nil {
		sandbox.Deps.Error(list_error.Error())
	}
	if list_error != nil {
		return cli.ExitFailure
	}
	return cli.ExitOk
}
