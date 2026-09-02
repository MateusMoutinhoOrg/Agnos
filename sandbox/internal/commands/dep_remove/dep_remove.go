package dep_remove

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depRemoveAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_remove"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"dep-remove"},
		Category:              "Core Commands",
		Args: []api.CliArg{
			{
				Id:          "dep",
				Description: "the dep to remove from the project",
				Examples: []string{
					config.ProjectName + " dep-remove embed",
				},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " dep-remove embed .",
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
					config.ProjectName + " dep-remove embed -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Removes an embedded dep from the project",
		LongDescription: "Removes every file that assets/deplist/<dep> installs into the\nproject, then calls build.",
		Examples: []string{
			config.ProjectName + " dep-remove embed",
			config.ProjectName + " dep-remove embed .",
		},
		Handler: CommandHander,
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	depArg := entries.GetArgById("dep")
	dep := depArg.Values[0].String()
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	remove_error := depRemoveAction.DepRemove(deps, path, dep)

	if !quietFlag.Exist && remove_error != nil {
		deps.Std.Error(remove_error.Error())
	}
	if remove_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
