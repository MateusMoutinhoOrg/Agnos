package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"build"},
		Category:              "Core Commands",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " build . ",
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
					config.ProjectName + " build -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			{
				Id:               "unsafe",
				ValidIdentifiers: []string{"--unsafe"},
				Description:      "Skips the verify schema gate before building",
				Examples: []string{
					config.ProjectName + " build --unsafe",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Build the project in a directory",
		LongDescription: "Builds the project in the given directory, compiling\nthe source code into the output artifacts. If no\npath is provided, the current directory is used.",
		Examples: []string{
			config.ProjectName + " build",
			config.ProjectName + " build .",
			config.ProjectName + " build ./my-project",
			config.ProjectName + " build -q",
		},
		Handler: CommandHander,
	}
}
func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	unsafeFlag := entries.GetFlagById("unsafe")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	if !unsafeFlag.Exist {
		if verify_error := verifyAction.Verify(deps, path); verify_error != nil {
			if !quietFlag.Exist {
				deps.Std.Error(verify_error.Error())
			}
			return api.ExitFailure
		}
	}

	build_error := buildAction.Build(deps, path)

	if !quietFlag.Exist && build_error != nil {
		deps.Std.Error(build_error.Error())
	}
	if build_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
