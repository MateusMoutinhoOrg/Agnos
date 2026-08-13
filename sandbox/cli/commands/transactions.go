package commands

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// Transactions runs the `transactions` command, listing every stored
// transaction or, when a category name follows, only that category's.
func Transactions(l *api.Lib) int {
	verb := l.Deps.VerbLib

	listed := []api.Transaction{}
	name, err := verb.GetNextStringArg()
	if err != nil {
		listed = l.ListTransactions()
	} else {
		stored, found := l.GetCategory(name)
		if !found {
			return Failure(l, config.CategoryNotFound, name)
		}
		listed = stored.ListTransactions()
	}

	if len(listed) == 0 {
		l.Deps.Printf("%s\n", config.NoTransactions)
		return api.ExitOk
	}
	for _, written := range listed {
		l.Deps.Printf("%s\n", written.String())
	}
	return api.ExitOk
}
