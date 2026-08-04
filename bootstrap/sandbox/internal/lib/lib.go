package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps"
)

// TestFuncFactory returns the closure that fills api.Lib.TestFunc, which
// exercises the embedded library reached through the Deps: it records one
// transaction, reads the tracker back, and prints it. The embedded library
// is never imported here — the adapter injects it as an agnosdeps.Lib
// struct, so calling it is just calling a function field.
func TestFuncFactory(l *api.Lib) func() {
	return func() {
		tracker := l.Deps.TrackerLib

		category, created := tracker.AddCategory("groceries")
		if !created {
			l.Deps.Println("groceries: could not create the category")
			return
		}
		category.AddSpend("weekly shopping", 8450)

		for _, transaction := range tracker.ListTransactions() {
			l.Deps.Println("transaction:", transaction.String())
		}
		l.Deps.Println("balance:", tracker.Balance())
	}
}

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.TestFunc = TestFuncFactory(&l)
	return l
}
