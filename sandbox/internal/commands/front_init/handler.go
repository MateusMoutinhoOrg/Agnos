package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	init_error := frontInitAction.FrontInit(deps, entries.Path)

	if init_error != nil {
		deps.Std.Error("%s\n", init_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
