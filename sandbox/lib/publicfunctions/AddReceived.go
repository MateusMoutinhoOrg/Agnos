package publicfunctions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// AddReceivedFactory fills api.Lib.AddReceived with a closure that records
// money entering the budget under an existing category. It reports false
// when the category is unknown.
func AddReceivedFactory(l *api.Lib) func(category string, description string, amount int64) (api.Transaction, bool) {
	return func(categoryName string, description string, amount int64) (api.Transaction, bool) {
		stored, ok := l.GetCategory(categoryName)
		if !ok {
			return api.Transaction{}, false
		}
		return stored.AddReceived(description, amount)
	}
}
