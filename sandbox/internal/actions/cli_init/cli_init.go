package cli_init

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CliInit installs the std and argv deps the CLI layer depends on and renders
// the "cli" asset group into the project, then runs build as a follow-up step.
func CliInit(deps *deps.Deps, path string) error {
	for _, dep := range []string{"std", "verb"} {
		if err := depInstallAction.DepInstall(deps, path, dep); err != nil {
			return err
		}
	}
	io := smartio.New(deps, path, config.ProjectName)
	if err := CliInitInternal(deps, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, path)
}
