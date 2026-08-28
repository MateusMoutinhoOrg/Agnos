package install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{

		ValidStartIdentifiers: []string{"install"},
		Category:              "Extensions",

		Args: []api.CliArg{
			api.CliArg{
				Id:          "item",
				Description: "the extension to install in the project",
				Examples: []string{
					sandbox.ProjectName + " install my-extension",
				},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
		},
		Flags: []api.Cliflag{
			api.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.ProjectName + " install my-extension -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			api.Cliflag{
				Id:               "path",
				ValidIdentifiers: []string{"--path", "-p"},
				Description:      "the dir of the project to install the extension into",
				Examples: []string{
					sandbox.ProjectName + " install my-extension -p ./my-project",
				},
				Defaults:         []string{"."},
				Type:             api.CliTypeString,
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

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	pathFlag := entries.GetFlagById("path")
	itemArg := entries.GetArgById("item")

	path := "."
	if pathFlag.Exist && len(pathFlag.Values) > 0 {
		path = pathFlag.Values[0].String()
	}

	install_error := sandbox.Install(api.InstallProps{
		Path: path,
		Item: itemArg.Values[0].String(),
	})

	if !quietFlag.Exist && install_error != nil {
		sandbox.Deps.Error(install_error.Error())
	}
	if install_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
