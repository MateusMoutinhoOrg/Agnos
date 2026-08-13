package commands

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/lib/store"
)

// Balance runs the `balance` command, printing the total across every
// category or, when a category name follows, that category's own balance.
func Balance(l *api.Lib) int {
	verb := l.Deps.VerbLib

	name, err := verb.GetNextStringArg()
	if err != nil {
		l.Deps.Printf("%s\n", store.Money(l.Balance()))
		return api.ExitOk
	}

	stored, found := l.GetCategory(name)
	if !found {
		return Failure(l, config.CategoryNotFound, name)
	}
	l.Deps.Printf("%s\n", store.Money(stored.Balance()))
	return api.ExitOk
}
