package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/category"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/store"
)

// AddCategoryFactory fills api.Lib.AddCategory with a closure that creates
// the named category in the injected database and returns it. Creation is
// idempotent: a name already taken makes the injected database report a key
// conflict, and the stored category is returned instead of a failure.
func AddCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		if name == "" {
			return api.Category{}, false
		}
		categories, ok := store.Categories(l.Deps)
		if !ok {
			return api.Category{}, false
		}

		record, err := categories.NewItem(map[string]any{
			store.NameField:      name,
			store.CreatedAtField: l.Deps.Now().Unix(),
		})
		if err == nil {
			return category.New(l.Deps, record), true
		}
		if err.Type != keepdeps.KeyConflict {
			return api.Category{}, false
		}

		stored, found := categories.FindByKey(store.NameField, name)
		if !found {
			return api.Category{}, false
		}
		return category.New(l.Deps, stored), true
	}
}

// GetCategoryFactory fills api.Lib.GetCategory with a closure that looks a
// stored category up by its unique name. It reports false when no category
// carries that name.
func GetCategoryFactory(l *api.Lib) func(name string) (api.Category, bool) {
	return func(name string) (api.Category, bool) {
		record, ok := store.FindCategory(l.Deps, name)
		if !ok {
			return api.Category{}, false
		}
		return category.New(l.Deps, record), true
	}
}

// ListCategoriesFactory fills api.Lib.ListCategories with a closure that
// reads every stored category back out of the injected database.
func ListCategoriesFactory(l *api.Lib) func() []api.Category {
	return func() []api.Category {
		categories, ok := store.Categories(l.Deps)
		if !ok {
			return nil
		}
		records, err := categories.ListAll()
		if err != nil {
			return nil
		}
		listed := make([]api.Category, 0, len(records))
		for _, record := range records {
			listed = append(listed, category.New(l.Deps, record))
		}
		return listed
	}
}

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

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.AddCategory = AddCategoryFactory(&l)
	l.GetCategory = GetCategoryFactory(&l)
	l.ListCategories = ListCategoriesFactory(&l)
	l.AddSpend = AddSpendFactory(&l)
	l.AddReceived = AddReceivedFactory(&l)
	l.ListTransactions = ListTransactionsFactory(&l)
	l.Balance = BalanceFactory(&l)
	return l
}
