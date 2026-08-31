package uninstall

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{

		ValidStartIdentifiers: []string{"uninstall"},
		Category:              "Extensions",

		Args: []cli.CliArg{
			cli.CliArg{
				Id:          "item",
				Description: "the extension to uninstall from the project",
				Examples: []string{
					sandbox.Config.ProjectName + " uninstall my-extension",
				},
				RequiredType:    cli.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
		},
		Flags: []cli.Cliflag{
			cli.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.Config.ProjectName + " uninstall my-extension -q",
				},
				Type:             cli.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			cli.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project to uninstall the extension from",
				Examples: []string{
					sandbox.Config.ProjectName + " uninstall my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             cli.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Uninstall an extension from the project",
		LongDescription: "Removes the given extension from the project, deleting the files\nand configuration it added. If no path is provided, the current\ndirectory is used.",
		Examples: []string{
			sandbox.Config.ProjectName + " uninstall my-extension",
			sandbox.Config.ProjectName + " uninstall my-extension -p ./my-project",
			sandbox.Config.ProjectName + " uninstall my-extension -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)

	quietFlag := entries.GetFlagById("quiet")
	pathFlag := entries.GetFlagById("path")
	itemArg := entries.GetArgById("item")

	path := "."
	if pathFlag.Exist && len(pathFlag.Values) > 0 {
		path = pathFlag.Values[0].String()
	}

	uninstall_error := sandbox.Core.Uninstall(core.UninstallProps{
		Path: path,
		Item: itemArg.Values[0].String(),
	})

	if !quietFlag.Exist && uninstall_error != nil {
		sandbox.Deps.Error(uninstall_error.Error())
	}
	if uninstall_error != nil {
		return cli.ExitFailure
	}
	return cli.ExitOk
}
