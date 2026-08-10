package agnosdeps

import "time"

// This package is the bootstrap library's *copy* of the embedded Agnos-Cli
// financial-tracker library's public api. The sandbox may not import the
// embedded library — that would be a third-party import — so it restates the
// shape it needs here, field for field. The adapter, which lives outside the
// sandbox, is what fills these structs from the real library.
//
// Copying is cheap precisely because the embedded library exposes structs
// of function fields instead of interfaces: an adapter assigns the real
// library's fields straight into the copy. Where a field hands back another
// api struct — Lib.AddCategory, Category.AddSpend — the adapter wraps it in a
// closure that copies the returned struct too.

// Transaction kinds, reported by Transaction.Kind. The values mirror the
// embedded library's constants, so they can be assigned across without
// mapping.
const (
	// Spend is money leaving the tracked budget.
	Spend = iota
	// Received is money entering the tracked budget.
	Received
)

// Transaction mirrors the embedded library's api.Transaction — one spend or
// received record. The embedded library's Deps field is deliberately not
// copied: the sandbox has no business reading the dependencies the embedded
// library was wired with.
type Transaction struct {
	// Id is the record's permanent identifier inside its category.
	Id int64
	// Category is the name of the category the transaction belongs to.
	Category string
	// Reference is the record's unique key inside its category.
	Reference string
	// Description is the human-readable note the transaction was recorded
	// with.
	Description string
	// Amount is the value of the transaction in the smallest currency unit,
	// always positive.
	Amount int64
	// Kind is Spend or Received.
	Kind int
	// OccurredAt is the moment the transaction was recorded.
	OccurredAt time.Time
	// SignedAmount returns Amount negated for a Spend and unchanged for a
	// Received.
	SignedAmount func() int64
	// Remove deletes the transaction from its category.
	Remove func() bool
	// String renders the transaction as one human-readable line.
	String func() string
}

// Category mirrors the embedded library's api.Category — one bucket
// transactions are tracked under.
type Category struct {
	// Id is the record's permanent identifier.
	Id int64
	// Name is the category's unique name.
	Name string
	// CreatedAt is the moment the category was created.
	CreatedAt time.Time
	// AddSpend records money leaving the budget under this category.
	AddSpend func(description string, amount int64) (Transaction, bool)
	// AddReceived records money entering the budget under this category.
	AddReceived func(description string, amount int64) (Transaction, bool)
	// ListTransactions returns every transaction stored under this category.
	ListTransactions func() []Transaction
	// Balance sums the signed amounts of this category's transactions.
	Balance func() int64
	// Remove deletes the category and every transaction under it.
	Remove func() bool
	// String renders the category as one human-readable line.
	String func() string
}

// Lib mirrors the embedded library's api.Lib — a financial tracker over a
// persisted schema database.
type Lib struct {
	// Sandboxmain is the embedded library's command-line interface, run over
	// an argument vector and returning a process exit code.
	Sandboxmain func(args []string) int
	// AddCategory creates the named category, or returns the stored one when
	// the name is already taken.
	AddCategory func(name string) (Category, bool)
	// GetCategory returns the stored category with the given name.
	GetCategory func(name string) (Category, bool)
	// ListCategories returns every stored category.
	ListCategories func() []Category
	// AddSpend records money leaving the budget under an existing category.
	AddSpend func(category string, description string, amount int64) (Transaction, bool)
	// AddReceived records money entering the budget under an existing
	// category.
	AddReceived func(category string, description string, amount int64) (Transaction, bool)
	// ListTransactions returns every transaction of every category.
	ListTransactions func() []Transaction
	// Balance sums the signed amounts of every stored transaction.
	Balance func() int64
}
