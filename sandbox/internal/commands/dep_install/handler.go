package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_install"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	install_error := depInstallAction.DepInstall(deps, entries.Path, entries.Dep)

	if install_error != nil {
		deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
