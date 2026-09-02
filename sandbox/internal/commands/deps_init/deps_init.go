package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_init"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"deps-init"},
		Category:              "Dependency System",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " deps-init . ",
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
					config.ProjectName + " deps-init -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Initializes the dependency-injection subsystem for the project",
		LongDescription: "Creates the sandbox/deps and adapters directories and calls build.\nRun this once before using dep-install.",
		Examples: []string{
			config.ProjectName + " deps-init",
			config.ProjectName + " deps-init .",
		},
		Handler: CommandHander,
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	init_error := depsInitAction.DepsInit(deps, path)

	if !quietFlag.Exist && init_error != nil {
		deps.Std.Error(init_error.Error())
	}
	if init_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
