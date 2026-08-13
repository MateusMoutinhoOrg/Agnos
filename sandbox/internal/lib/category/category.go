package category

import (
	"strconv"
	"strings"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/lib/store"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/lib/transaction"
)

// referenceAttempts bounds how many sequence numbers a new transaction may
// try before giving up. A removal leaves a gap in the sequence, so the first
// candidate can collide with a record already stored under that reference.
const referenceAttempts = 64

// AddSpendFactory returns the closure that fills api.Category.AddSpend,
// recording money leaving the budget under this category.
func AddSpendFactory(c *api.Category) func(description string, amount int64) (api.Transaction, bool) {
	return func(description string, amount int64) (api.Transaction, bool) {
		return add(c, api.Spend, description, amount)
	}
}

// AddReceivedFactory returns the closure that fills
// api.Category.AddReceived, recording money entering the budget under this
// category.
func AddReceivedFactory(c *api.Category) func(description string, amount int64) (api.Transaction, bool) {
	return func(description string, amount int64) (api.Transaction, bool) {
		return add(c, api.Received, description, amount)
	}
}

// ListTransactionsFactory returns the closure that fills
// api.Category.ListTransactions, reading the category's nested collection
// back out of the injected database on every call.
func ListTransactionsFactory(c *api.Category) func() []api.Transaction {
	return func() []api.Transaction {
		record, ok := store.FindCategory(c.Deps, c.Name)
		if !ok {
			return nil
		}
		records := record.ListAll(store.TransactionsField)
		transactions := make([]api.Transaction, 0, len(records))
		for _, stored := range records {
			transactions = append(transactions, transaction.New(c.Deps, c.Name, stored))
		}
		return transactions
	}
}

// BalanceFactory returns the closure that fills api.Category.Balance,
// summing the signed amounts of the category's transactions.
func BalanceFactory(c *api.Category) func() int64 {
	return func() int64 {
		balance := int64(0)
		for _, t := range c.ListTransactions() {
			balance += t.SignedAmount()
		}
		return balance
	}
}

// RemoveFactory returns the closure that fills api.Category.Remove, deleting
// the category and — because Keep removes a record's nested collections with
// it — every transaction stored under it.
func RemoveFactory(c *api.Category) func() bool {
	return func() bool {
		record, ok := store.FindCategory(c.Deps, c.Name)
		if !ok {
			return false
		}
		return record.Remove() == nil
	}
}

// StringFactory returns the closure that fills api.Category.String,
// rendering the category as one line: its name, its live balance, and how
// many transactions it holds.
func StringFactory(c *api.Category) func() string {
	return func() string {
		transactions := c.ListTransactions()
		balance := int64(0)
		for _, t := range transactions {
			balance += t.SignedAmount()
		}
		fields := []string{
			c.Name,
			store.Money(balance),
			strconv.Itoa(len(transactions)) + " transactions",
		}
		return strings.Join(fields, "  ")
	}
}

// add writes one transaction of the given kind into the category's nested
// collection and returns it. It is the unexported helper both AddSpend and
// AddReceived close over: the two differ only by the kind they record.
//
// The description is not unique, so it travels inside the record's unique
// reference together with a sequence number; a reference already taken makes
// Keep report a key conflict, and the next sequence number is tried.
func add(c *api.Category, kind int, description string, amount int64) (api.Transaction, bool) {
	if amount <= 0 {
		return api.Transaction{}, false
	}
	record, ok := store.FindCategory(c.Deps, c.Name)
	if !ok {
		return api.Transaction{}, false
	}

	occurredAt := c.Deps.Now()
	sequence := int64(len(record.ListAll(store.TransactionsField))) + 1
	for attempt := 0; attempt < referenceAttempts; attempt++ {
		created, err := record.NewSubItem(store.TransactionsField, map[string]any{
			store.ReferenceField:  store.Reference(sequence+int64(attempt), description),
			store.AmountField:     amount,
			store.KindField:       int64(kind),
			store.OccurredAtField: occurredAt.Unix(),
		})
		if err == nil {
			return transaction.New(c.Deps, c.Name, created), true
		}
		if err.Type != keepdeps.KeyConflict {
			return api.Transaction{}, false
		}
	}
	return api.Transaction{}, false
}

// New builds an api.Category from a stored record, propagating the library's
// Deps into it, and runs every category factory over it to fill its function
// fields. Adding a function field to api.Category means adding its factory
// call here.
func New(d deps.Deps, record keepdeps.SchemaItem) api.Category {
	c := api.Category{
		Deps:      d,
		Id:        record.Id,
		Name:      store.Text(record, store.NameField),
		CreatedAt: time.Unix(store.Number(record, store.CreatedAtField), 0),
	}
	c.AddSpend = AddSpendFactory(&c)
	c.AddReceived = AddReceivedFactory(&c)
	c.ListTransactions = ListTransactionsFactory(&c)
	c.Balance = BalanceFactory(&c)
	c.Remove = RemoveFactory(&c)
	c.String = StringFactory(&c)
	return c
}
