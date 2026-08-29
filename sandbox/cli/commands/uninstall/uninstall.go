package uninstall

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func NewCommand(sandbox *lib.SandBox) lib.CliCommand {
	return lib.CliCommand{

		ValidStartIdentifiers: []string{"uninstall"},
		Category:              "Extensions",

		Args: []lib.CliArg{
			lib.CliArg{
				Id:          "item",
				Description: "the extension to uninstall from the project",
				Examples: []string{
					sandbox.ProjectName + " uninstall my-extension",
				},
				RequiredType:    lib.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
		},
		Flags: []lib.Cliflag{
			lib.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.ProjectName + " uninstall my-extension -q",
				},
				Type:             lib.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			lib.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project to uninstall the extension from",
				Examples: []string{
					sandbox.ProjectName + " uninstall my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             lib.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Uninstall an extension from the project",
		LongDescription: "Removes the given extension from the project, deleting the files\nand configuration it added. If no path is provided, the current\ndirectory is used.",
		Examples: []string{
			sandbox.ProjectName + " uninstall my-extension",
			sandbox.ProjectName + " uninstall my-extension -p ./my-project",
			sandbox.ProjectName + " uninstall my-extension -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *lib.SandBox, entries lib.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	pathFlag := entries.GetFlagById("path")
	itemArg := entries.GetArgById("item")

	path := "."
	if pathFlag.Exist && len(pathFlag.Values) > 0 {
		path = pathFlag.Values[0].String()
	}

	uninstall_error := sandbox.Uninstall(lib.UninstallProps{
		Path: path,
		Item: itemArg.Values[0].String(),
	})

	if !quietFlag.Exist && uninstall_error != nil {
		sandbox.Deps.Error(uninstall_error.Error())
	}
	if uninstall_error != nil {
		return lib.ExitFailure
	}
	return lib.ExitOk
}
