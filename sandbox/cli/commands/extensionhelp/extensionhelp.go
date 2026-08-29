package extensionhelp

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{

		ValidStartIdentifiers: []string{"extension-help"},
		Category:              "Extensions",

		Args: []cli.CliArg{
			cli.CliArg{
				Id:          "extension",
				Description: "the extension to show the help of",
				Examples: []string{
					sandbox.Config.ProjectName + " extension-help my-extension",
				},
				RequiredType:    cli.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
		},
		Flags: []cli.Cliflag{
			cli.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project the extension belongs to",
				Examples: []string{
					sandbox.Config.ProjectName + " extension-help my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             cli.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Show the help of an extension",
		LongDescription: "Prints the help of the given extension: what it does, what it\ninstalls, and how it is used. If no path is provided, the current\ndirectory is used.",
		Examples: []string{
			sandbox.Config.ProjectName + " extension-help my-extension",
			sandbox.Config.ProjectName + " extension-help my-extension -p ./my-project",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)

	pathFlag := entries.GetFlagById("path")
	extensionArg := entries.GetArgById("extension")

	path := "."
	if pathFlag.Exist && len(pathFlag.Values) > 0 {
		path = pathFlag.Values[0].String()
	}

	help_error := sandbox.Core.ExtensionHelp(core.ExtensionHelpProps{
		Path:      path,
		Extension: extensionArg.Values[0].String(),
	})

	if help_error != nil {
		sandbox.Deps.Error(help_error.Error())
		return cli.ExitFailure
	}
	return cli.ExitOk
}
