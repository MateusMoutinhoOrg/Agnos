package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/lib/publicfunctions"
)

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.Sandboxmain = publicfunctions.SandboxmainFactory(&l)
	l.AddCategory = publicfunctions.AddCategoryFactory(&l)
	l.GetCategory = publicfunctions.GetCategoryFactory(&l)
	l.ListCategories = publicfunctions.ListCategoriesFactory(&l)
	l.AddSpend = publicfunctions.AddSpendFactory(&l)
	l.AddReceived = publicfunctions.AddReceivedFactory(&l)
	l.ListTransactions = publicfunctions.ListTransactionsFactory(&l)
	l.Balance = publicfunctions.BalanceFactory(&l)
	return l
}
