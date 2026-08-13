package api

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// Transaction kinds, reported by Transaction.Kind. The kind is what carries
// the direction of the money: Amount itself is always positive, and
// Transaction.SignedAmount applies the sign.
const (
	// Spend is money leaving the tracked budget.
	Spend = iota
	// Received is money entering the tracked budget.
	Received
)

// Transaction is a single spend or received record handed back by the
// library, already persisted in the injected database. Deps is the injected
// dependency set the transaction was built with; the plain data fields carry
// the record itself; the function fields are filled by factories in
// sandbox/lib/transaction, each closing over this struct so it reads
// Deps and the record's data at call time.
type Transaction struct {
	// Deps is the dependency set injected into the library, carried here so
	// the transaction's factories can reach it.
	Deps deps.Deps
	// Id is the record's permanent identifier inside its category.
	Id int64
	// Category is the name of the category the transaction belongs to.
	Category string
	// Reference is the record's unique key inside its category, composed by
	// the library from a sequence number and the description.
	Reference string
	// Description is the human-readable note the transaction was recorded
	// with.
	Description string
	// Amount is the value of the transaction in the smallest currency unit
	// (cents), always positive.
	Amount int64
	// Kind is Spend or Received.
	Kind int
	// OccurredAt is the moment the transaction was recorded, stamped from the
	// injected clock.
	OccurredAt time.Time
	// SignedAmount returns Amount negated for a Spend and unchanged for a
	// Received, so transactions can be summed directly.
	SignedAmount func() int64
	// Remove deletes the transaction from its category, reporting whether the
	// record was found and deleted.
	Remove func() bool
	// String renders the transaction as one human-readable line.
	String func() string
}

// Category is one bucket transactions are tracked under — "groceries",
// "salary" — already persisted in the injected database. Its function fields
// are filled by factories in sandbox/lib/category, each closing over
// this struct so every call re-reads the stored record through Deps.
type Category struct {
	// Deps is the dependency set injected into the library, carried here so
	// the category's factories can reach it.
	Deps deps.Deps
	// Id is the record's permanent identifier.
	Id int64
	// Name is the category's unique name, as passed to Lib.AddCategory.
	Name string
	// CreatedAt is the moment the category was created, stamped from the
	// injected clock.
	CreatedAt time.Time
	// AddSpend records money leaving the budget under this category. amount
	// is in the smallest currency unit (cents) and must be positive; the bool
	// is false when the amount is not positive or the record could not be
	// written.
	AddSpend func(description string, amount int64) (Transaction, bool)
	// AddReceived records money entering the budget under this category,
	// with the same rules as AddSpend.
	AddReceived func(description string, amount int64) (Transaction, bool)
	// ListTransactions returns every transaction stored under this category,
	// oldest first.
	ListTransactions func() []Transaction
	// Balance sums the signed amounts of this category's transactions.
	Balance func() int64
	// Remove deletes the category and every transaction stored under it,
	// reporting whether the record was found and deleted.
	Remove func() bool
	// String renders the category as one human-readable line.
	String func() string
}

// Exit codes reported by Lib.Sandboxmain, and through it by the process that
// calls it.
const (
	// ExitOk reports that the requested command ran to completion.
	ExitOk = 0
	// ExitUsage reports that the command line itself was wrong — an unknown
	// command, a missing operand, an unparsable amount — and the usage
	// screen was printed.
	ExitUsage = 1
	// ExitFailure reports that a well-formed command could not be carried
	// out, because a record was missing or could not be written.
	ExitFailure = 2
)

// Lib is the entry point handed back by lib.New. It is a financial tracker:
// categories hold spend and received transactions, and every record is
// persisted through the schema database injected as Deps.KeepLib. It is
// exposed as a struct of function fields: lib.New stores the injected deps in
// Deps and then runs the factories in sandbox/lib, each of which
// fills one function field with a closure over this struct.
//
// Because it is a struct and not an interface, a consumer that itself uses
// this pattern can copy the shape of Lib into its own deps contract and
// receive the whole library as a single injected field.
type Lib struct {
	// Deps is the dependency set injected by lib.New, carried here so every
	// factory-built function field can reach it.
	Deps deps.Deps
	// Sandboxmain is the command-line interface: the whole program, run
	// inside the sandbox. It reads the actions and flags of args through the
	// injected Deps.VerbLib parser, calls the library functions below, prints
	// every result and error through Deps.Printf, and returns the process
	// exit code — ExitOk, ExitUsage, or ExitFailure. The caller in cmd/main
	// does nothing but hand it the argument vector and exit with what it
	// returns.
	//
	// args must be the same argument vector the adapter wired Deps.VerbLib
	// over: the parser owns the reading, and args is what Sandboxmain checks
	// for an empty command line. The standard adapter and cmd/main both take
	// it from os.Args[1:], so they agree by construction.
	Sandboxmain func(args []string) int
	// AddCategory creates the category with the given name and returns it.
	// Creation is idempotent: when the name is already taken the stored
	// category is returned instead. The bool is false when the name is empty
	// or the record could not be written.
	AddCategory func(name string) (Category, bool)
	// GetCategory returns the stored category with the given name. The bool
	// is false when no category carries that name.
	GetCategory func(name string) (Category, bool)
	// ListCategories returns every stored category, oldest first.
	ListCategories func() []Category
	// AddSpend records money leaving the budget under an existing category.
	// The bool is false when the category is unknown, the amount is not
	// positive, or the record could not be written.
	AddSpend func(category string, description string, amount int64) (Transaction, bool)
	// AddReceived records money entering the budget under an existing
	// category, with the same rules as AddSpend.
	AddReceived func(category string, description string, amount int64) (Transaction, bool)
	// ListTransactions returns every transaction of every category, grouped
	// by category in creation order.
	ListTransactions func() []Transaction
	// Balance sums the signed amounts of every stored transaction.
	Balance func() int64
}
