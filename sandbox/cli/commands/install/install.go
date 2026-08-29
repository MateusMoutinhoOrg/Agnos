package install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func NewCommand(sandbox *lib.SandBox) lib.CliCommand {
	return lib.CliCommand{

		ValidStartIdentifiers: []string{"install"},
		Category:              "Extensions",

		Args: []lib.CliArg{
			lib.CliArg{
				Id:          "item",
				Description: "the extension to install in the project",
				Examples: []string{
					sandbox.ProjectName + " install my-extension",
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
					sandbox.ProjectName + " install my-extension -q",
				},
				Type:             lib.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			lib.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project to install the extension into",
				Examples: []string{
					sandbox.ProjectName + " install my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             lib.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Install an extension in the project",
		LongDescription: "Installs the given extension in the project, adding the files\nand configuration it needs. If no path is provided, the current\ndirectory is used.",
		Examples: []string{
			sandbox.ProjectName + " install my-extension",
			sandbox.ProjectName + " install my-extension -p ./my-project",
			sandbox.ProjectName + " install my-extension -q",
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

	install_error := sandbox.Install(lib.InstallProps{
		Path: path,
		Item: itemArg.Values[0].String(),
	})

	if !quietFlag.Exist && install_error != nil {
		sandbox.Deps.Error(install_error.Error())
	}
	if install_error != nil {
		return lib.ExitFailure
	}
	return lib.ExitOk
}
