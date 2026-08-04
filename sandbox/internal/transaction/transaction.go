package transaction

import (
	"strings"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/store"
)

// SignedAmountFactory returns the closure that fills
// api.Transaction.SignedAmount, negating the stored amount for a spend so a
// list of transactions can be summed directly.
func SignedAmountFactory(t *api.Transaction) func() int64 {
	return func() int64 {
		if t.Kind == api.Spend {
			return -t.Amount
		}
		return t.Amount
	}
}

// RemoveFactory returns the closure that fills api.Transaction.Remove,
// deleting the record from its category's nested collection. It reports
// false when the category or the record is no longer stored.
func RemoveFactory(t *api.Transaction) func() bool {
	return func() bool {
		category, ok := store.FindCategory(t.Deps, t.Category)
		if !ok {
			return false
		}
		for _, record := range category.ListAll(store.TransactionsField) {
			if store.Text(record, store.ReferenceField) != t.Reference {
				continue
			}
			return record.Remove() == nil
		}
		return false
	}
}

// StringFactory returns the closure that fills api.Transaction.String,
// rendering the record as one line: when it happened, its direction, its
// category, its signed amount, and its description.
func StringFactory(t *api.Transaction) func() string {
	return func() string {
		fields := []string{
			t.OccurredAt.Format("2006-01-02 15:04"),
			kindName(t.Kind),
			t.Category,
			store.Money(t.SignedAmount()),
			t.Description,
		}
		return strings.Join(fields, "  ")
	}
}

// kindName renders a transaction kind for StringFactory. It is an
// unexported helper, not a field factory: no api field holds it.
func kindName(kind int) string {
	if kind == api.Spend {
		return "spend"
	}
	return "received"
}

// New builds an api.Transaction from a stored record, propagating the
// library's Deps into it, and runs every transaction factory over it to fill
// its function fields. categoryName is the category the record is nested
// under, which the record itself does not carry. Adding a function field to
// api.Transaction means adding its factory call here.
func New(d deps.Deps, categoryName string, record keepdeps.SchemaItem) api.Transaction {
	reference := store.Text(record, store.ReferenceField)
	t := api.Transaction{
		Deps:        d,
		Id:          record.Id,
		Category:    categoryName,
		Reference:   reference,
		Description: store.Description(reference),
		Amount:      store.Number(record, store.AmountField),
		Kind:        int(store.Number(record, store.KindField)),
		OccurredAt:  time.Unix(store.Number(record, store.OccurredAtField), 0),
	}
	t.SignedAmount = SignedAmountFactory(&t)
	t.Remove = RemoveFactory(&t)
	t.String = StringFactory(&t)
	return t
}
