package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"dep-install"},
		Category:              "Dependencies",
		Args: []api.CliArg{
			{
				Id:          "dep",
				Description: "the dep to install from assets/deplist",
				Examples: []string{
					config.ProjectName + " dep-install embed",
				},
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 1,
				RequiredMaxSize: 1,
			},
			{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					config.ProjectName + " dep-install embed .",
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
					config.ProjectName + " dep-install embed -q",
				},
				Type:             api.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},
		Description:     "Installs an embedded dep into the project",
		LongDescription: "Renders every file under assets/deplist/<dep> into the project\nat the path it holds inside that dep, then calls build.",
		Examples: []string{
			config.ProjectName + " dep-install embed",
			config.ProjectName + " dep-install embed .",
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

	install_error := depInstallAction.DepInstall(deps, path, dep)

	if !quietFlag.Exist && install_error != nil {
		deps.Std.Error(install_error.Error())
	}
	if install_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
