package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{

		ValidStartIdentifiers: []string{"start"},

		Args: []api.CliArg{
			api.CliArg{
				Id:          "path",
				Description: "the dir to start the project",
				Examples: []string{
					sandbox.ProjectName + " start . ",
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
					sandbox.ProjectName + " start -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},

		Description: "Starts the agnos cli",
		Examples: []string{
			sandbox.ProjectName + " start",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	start_error := sandbox.Start(path)

	if !quietFlag.Exist && start_error != nil {
		sandbox.Deps.Printf(start_error.Error())
	}
	if start_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
