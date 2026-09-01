package remove_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	removeDepsAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"remove-deps"},
		Category:              "Core Commands",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " remove-deps . ",
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
					config.ProjectName + " remove-deps -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Removes dependencies for the project",
		LongDescription: "Removes sandbox/deps and adapters directories and calls build.",
		Examples: []string{
			config.ProjectName + " remove-deps",
			config.ProjectName + " remove-deps .",
		},
		Handler: CommandHander,
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	io := smartio.New(deps, path, config.ProjectName)
	
	build_error := removeDepsAction.RemoveDeps(deps, io, path)

	if build_error == nil {
		build_error = io.Persist()
	}

	if build_error == nil {
		build_error = buildAction.Build(deps, io, path)
	}

	if build_error == nil {
		build_error = io.Persist()
	}

	if !quietFlag.Exist && build_error != nil {
		deps.Error(build_error.Error())
	}
	if build_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
