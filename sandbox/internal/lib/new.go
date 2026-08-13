package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.Sandboxmain = SandboxmainFactory(&l)
	l.AddCategory = AddCategoryFactory(&l)
	l.GetCategory = GetCategoryFactory(&l)
	l.ListCategories = ListCategoriesFactory(&l)
	l.AddSpend = AddSpendFactory(&l)
	l.AddReceived = AddReceivedFactory(&l)
	l.ListTransactions = ListTransactionsFactory(&l)
	l.Balance = BalanceFactory(&l)
	return l
}
