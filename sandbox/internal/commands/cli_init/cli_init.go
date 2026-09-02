package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	cliInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/cli_init"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"cli-init"},
		Category:              "Cli System",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " cli-init . ",
				},
				Defaults:        []string{"."},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []api.Cliflag{
			{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					config.ProjectName + " cli-init -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Initializes the CLI layer for the project",
		LongDescription: "Installs the argv dep the CLI layer depends on, renders the\n\"cli\" asset group into the project, and calls build.",
		Examples: []string{
			config.ProjectName + " cli-init",
			config.ProjectName + " cli-init .",
		},
		Handler: CommandHander,
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	init_error := cliInitAction.CliInit(deps, path)

	if !quietFlag.Exist && init_error != nil {
		deps.Std.Error(init_error.Error())
	}
	if init_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
