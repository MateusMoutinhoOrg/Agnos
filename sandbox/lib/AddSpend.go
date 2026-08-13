package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// AddSpendFactory fills api.Lib.AddSpend with a closure that records money
// leaving the budget under an existing category. It reports false when the
// category is unknown.
func AddSpendFactory(l *api.Lib) func(category string, description string, amount int64) (api.Transaction, bool) {
	return func(categoryName string, description string, amount int64) (api.Transaction, bool) {
		stored, ok := l.GetCategory(categoryName)
		if !ok {
			return api.Transaction{}, false
		}
		return stored.AddSpend(description, amount)
	}
}
