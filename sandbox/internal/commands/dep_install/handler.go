package dep_install

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	depInstallAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/dep_install"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	install_error := depInstallAction.DepInstall(deps, api.DepInstallProps{
		Path:    entries.Path,
		Dep:     entries.Dep,
		Adapter: entries.Adapter,
	})

	if install_error != nil {
		deps.Std.Error("%s\n", install_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
