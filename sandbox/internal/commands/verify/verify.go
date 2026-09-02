package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"verify"},
		Category:              "Core Commands",
		Args: []api.CliArg{
			{
				Id:          "path",
				Description: "the dir holding the project to verify",
				Examples: []string{
					config.ProjectName + " verify . ",
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
					config.ProjectName + " verify -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Checks the project keeps the sandbox/adapter schema",
		LongDescription: "Verifies the structural rules the harness depends on: sandbox/ imports\nstay inside sandbox/, sandbox/ holds only api, binds, deps, internal and\nnew.go, sandbox/api and sandbox/deps import nothing external, every\nsandbox/binds file mirrors a sandbox/api file and declares only functions,\nand adapters/ holds only availables and libs. `agnos build` runs this as a\ngate unless --unsafe is passed.",
		Examples: []string{
			config.ProjectName + " verify",
			config.ProjectName + " verify .",
		},
		Handler: CommandHander,
	}
}

func CommandHander(deps *deps.Deps, entries api.CliEntrys) int {
	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	verify_error := verifyAction.Verify(deps, path)

	if !quietFlag.Exist && verify_error != nil {
		deps.Std.Error(verify_error.Error())
	}
	if verify_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
