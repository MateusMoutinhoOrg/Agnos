package dep_list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depListAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_list"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"dep-list"},
		Category:              "Dependencies",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " dep-list .",
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
					config.ProjectName + " dep-list -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Lists the embedded deps available to install",
		LongDescription: "Lists the name of every dep under assets/deplist that dep-install\ncan render into a project.",
		Examples: []string{
			config.ProjectName + " dep-list",
		},
		Handler: CommandHander,
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	deplist, list_error := depListAction.DepList(deps, path)

	if !quietFlag.Exist && list_error != nil {
		deps.Std.Error(list_error.Error())
	}
	if list_error != nil {
		return api.ExitFailure
	}

	for _, dep := range deplist {
		deps.Std.Printf("%s\n", dep)
	}
	return api.ExitOk
}
