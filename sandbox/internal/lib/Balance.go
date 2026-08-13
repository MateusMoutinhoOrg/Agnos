package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// BalanceFactory fills api.Lib.Balance with a closure that sums the signed
// amounts of every stored transaction — received money counts up, spent
// money counts down.
func BalanceFactory(l *api.Lib) func() int64 {
	return func() int64 {
		balance := int64(0)
		for _, t := range l.ListTransactions() {
			balance += t.SignedAmount()
		}
		return balance
	}
}
