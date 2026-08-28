package extensionhelp

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{

		ValidStartIdentifiers: []string{"extension-help"},
		Category:              "Extensions",

		Args: []api.CliArg{
			api.CliArg{
				Id:          "extension",
				Description: "the extension to show the help of",
				Examples: []string{
					sandbox.ProjectName + " extension-help my-extension",
				},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
		},
		Flags: []api.Cliflag{
			api.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project the extension belongs to",
				Examples: []string{
					sandbox.ProjectName + " extension-help my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             api.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Show the help of an extension",
		LongDescription: "Prints the help of the given extension: what it does, what it\ninstalls, and how it is used. If no path is provided, the current\ndirectory is used.",
		Examples: []string{
			sandbox.ProjectName + " extension-help my-extension",
			sandbox.ProjectName + " extension-help my-extension -p ./my-project",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {

	pathFlag := entries.GetFlagById("path")
	extensionArg := entries.GetArgById("extension")

	path := "."
	if pathFlag.Exist && len(pathFlag.Values) > 0 {
		path = pathFlag.Values[0].String()
	}

	help_error := sandbox.ExtensionHelp(api.ExtensionHelpProps{
		Path:      path,
		Extension: extensionArg.Values[0].String(),
	})

	if help_error != nil {
		sandbox.Deps.Error(help_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
