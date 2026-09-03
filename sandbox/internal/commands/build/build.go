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
		Args:                  []api.CliArg{},
		Flags: []api.Cliflag{
			{
				Id:               "path",
				ValidIdentifiers: []string{"--path"},
				Description:      "the dir holding the project (defaults to the current directory)",
				Examples: []string{
					config.ProjectName + " build --path ./my-project",
				},
				Type:             api.CliTypeString,
				Defaults:         []string{"."},
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},

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
			config.ProjectName + " build --path ./my-project",
			config.ProjectName + " build -q",
		},
		Handler: func(entries api.CliEntrys) int { return CommandHander(deps, entries) },
	}
}
func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	unsafeFlag := entries.GetFlagById("unsafe")
	pathFlag := entries.GetFlagById("path")
	path := pathFlag.Values[0].String()

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
