package standard

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps/agnosdeps"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
	agnostypes "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// StandardAdapter fills the bootstrap deps.Deps using the Go standard library
// for output and the embedded Agnos financial-tracker library — wired with
// its own standard adapter — for storage. Only this file, which lives outside
// the bootstrap sandbox, may import the embedded library.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	// trackerBasePath is where the embedded tracker library persists its
	// categories and transactions.
	trackerBasePath string
}

// PrintlnFactory returns the closure that fills deps.Deps.Println, writing a
// line to standard output.
func PrintlnFactory(s *StandardAdapter) func(a ...any) {
	return func(a ...any) {
		fmt.Println(a...)
	}
}

// fromTransaction converts one transaction the embedded library handed back
// into the sandbox's copy. Its function fields cross over unwrapped: none of
// them returns another api struct.
func fromTransaction(transaction agnostypes.Transaction) agnosdeps.Transaction {
	return agnosdeps.Transaction{
		Id:           transaction.Id,
		Category:     transaction.Category,
		Reference:    transaction.Reference,
		Description:  transaction.Description,
		Amount:       transaction.Amount,
		Kind:         transaction.Kind,
		OccurredAt:   transaction.OccurredAt,
		SignedAmount: transaction.SignedAmount,
		Remove:       transaction.Remove,
		String:       transaction.String,
	}
}

// fromTransactions converts a slice of transactions into the sandbox's copy.
func fromTransactions(transactions []agnostypes.Transaction) []agnosdeps.Transaction {
	converted := make([]agnosdeps.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		converted = append(converted, fromTransaction(transaction))
	}
	return converted
}

// fromCategory converts one category the embedded library handed back into
// the sandbox's copy, wrapping every field that returns another api struct
// so nothing of the embedded library reaches the sandbox.
func fromCategory(category agnostypes.Category) agnosdeps.Category {
	return agnosdeps.Category{
		Id:        category.Id,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		AddSpend: func(description string, amount int64) (agnosdeps.Transaction, bool) {
			transaction, ok := category.AddSpend(description, amount)
			return fromTransaction(transaction), ok
		},
		AddReceived: func(description string, amount int64) (agnosdeps.Transaction, bool) {
			transaction, ok := category.AddReceived(description, amount)
			return fromTransaction(transaction), ok
		},
		ListTransactions: func() []agnosdeps.Transaction {
			return fromTransactions(category.ListTransactions())
		},
		Balance: category.Balance,
		Remove:  category.Remove,
		String:  category.String,
	}
}

// fromCategories converts a slice of categories into the sandbox's copy.
func fromCategories(categories []agnostypes.Category) []agnosdeps.Category {
	converted := make([]agnosdeps.Category, 0, len(categories))
	for _, category := range categories {
		converted = append(converted, fromCategory(category))
	}
	return converted
}

// TrackerLibFactory returns the value that fills deps.Deps.TrackerLib: the
// embedded financial-tracker library, initialized with its own standard
// adapter over the adapter's base path, copied field by field onto the
// sandbox's local agnosdeps.Lib. It returns a value rather than a closure
// because the field is a struct — see the Factories specification.
func TrackerLibFactory(s *StandardAdapter) agnosdeps.Lib {
	inner := agnoslib.New(agnosadapter.New(s.trackerBasePath, "."))
	return agnosdeps.Lib{
		Sandboxmain: inner.Sandboxmain,
		AddCategory: func(name string) (agnosdeps.Category, bool) {
			category, ok := inner.AddCategory(name)
			return fromCategory(category), ok
		},
		GetCategory: func(name string) (agnosdeps.Category, bool) {
			category, ok := inner.GetCategory(name)
			return fromCategory(category), ok
		},
		ListCategories: func() []agnosdeps.Category {
			return fromCategories(inner.ListCategories())
		},
		AddSpend: func(category string, description string, amount int64) (agnosdeps.Transaction, bool) {
			transaction, ok := inner.AddSpend(category, description, amount)
			return fromTransaction(transaction), ok
		},
		AddReceived: func(category string, description string, amount int64) (agnosdeps.Transaction, bool) {
			transaction, ok := inner.AddReceived(category, description, amount)
			return fromTransaction(transaction), ok
		},
		ListTransactions: func() []agnosdeps.Transaction {
			return fromTransactions(inner.ListTransactions())
		},
		Balance: inner.Balance,
	}
}

// New creates a deps.Deps backed by the bootstrap standard adapter, ready for
// lib.New. The embedded tracker library persists its records under the
// provided trackerBasePath. It builds the adapter instance and runs every
// field factory over it, assigning each return value into the matching field.
// Adding a field to deps.Deps means adding its factory call here.
func New(trackerBasePath string) deps.Deps {
	adapter := &StandardAdapter{trackerBasePath: trackerBasePath}
	adapter.Deps.Println = PrintlnFactory(adapter)
	adapter.Deps.TrackerLib = TrackerLibFactory(adapter)
	return adapter.Deps
}
