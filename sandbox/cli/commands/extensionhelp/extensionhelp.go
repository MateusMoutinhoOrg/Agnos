package extensionhelp

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func NewCommand(sandbox *lib.SandBox) lib.CliCommand {
	return lib.CliCommand{

		ValidStartIdentifiers: []string{"extension-help"},
		Category:              "Extensions",

		Args: []lib.CliArg{
			lib.CliArg{
				Id:          "extension",
				Description: "the extension to show the help of",
				Examples: []string{
					sandbox.ProjectName + " extension-help my-extension",
				},
				RequiredType:    lib.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
		},
		Flags: []lib.Cliflag{
			lib.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project the extension belongs to",
				Examples: []string{
					sandbox.ProjectName + " extension-help my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             lib.CliTypeString,
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

func CommandHandler(sandbox *lib.SandBox, entries lib.CliEntrys) int {

	pathFlag := entries.GetFlagById("path")
	extensionArg := entries.GetArgById("extension")

	path := "."
	if pathFlag.Exist && len(pathFlag.Values) > 0 {
		path = pathFlag.Values[0].String()
	}

	help_error := sandbox.ExtensionHelp(lib.ExtensionHelpProps{
		Path:      path,
		Extension: extensionArg.Values[0].String(),
	})

	if help_error != nil {
		sandbox.Deps.Error(help_error.Error())
		return lib.ExitFailure
	}
	return lib.ExitOk
}
