package publicfunctions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// ListTransactionsFactory fills api.Lib.ListTransactions with a closure that
// walks every stored category and collects its transactions.
func ListTransactionsFactory(l *api.Lib) func() []api.Transaction {
	return func() []api.Transaction {
		transactions := []api.Transaction{}
		for _, stored := range l.ListCategories() {
			transactions = append(transactions, stored.ListTransactions()...)
		}
		return transactions
	}
}
