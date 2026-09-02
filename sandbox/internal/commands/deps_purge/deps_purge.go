package deps_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depsPurgeAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/deps_purge"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"deps-purge"},
		Category:              "Dependency System",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " deps-purge . ",
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
					config.ProjectName + " deps-purge -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Removes the dependency-injection subsystem from the project",
		LongDescription: "Removes the sandbox/deps and adapters directories and calls build.",
		Examples: []string{
			config.ProjectName + " deps-purge",
			config.ProjectName + " deps-purge .",
		},
		Handler: func(entries api.CliEntrys) int { return CommandHander(deps, entries) },
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	purge_error := depsPurgeAction.DepsPurge(deps, path)

	if !quietFlag.Exist && purge_error != nil {
		deps.Std.Error(purge_error.Error())
	}
	if purge_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
