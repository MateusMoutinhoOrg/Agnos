package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_install"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	install_error := depInstallAction.DepInstall(deps, entries.Path, entries.Dep)

	if !entries.Quiet && install_error != nil {
		deps.Std.Error(install_error.Error())
	}
	if install_error != nil {
		return api.ExitFailure
	}
	return api.ExitOk
}
