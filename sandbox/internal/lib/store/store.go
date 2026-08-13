package store

// Shared database layout of the financial tracker: the schema every
// category and transaction is persisted under, plus the helpers the object
// packages use to reach it. Everything here is expressed against the
// embedded Keep library injected as Deps.KeepLib, so the sandbox stays
// closed — no adapter, no third-party module, no OS-bound package.
//
// This package declares no types; the object packages
// (sandbox/internal/lib/category, sandbox/internal/lib/transaction) and
// sandbox/internal/lib call into it.

import (
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"
)

const (
	// DatabasePath is the prefix every key the tracker writes lives under.
	DatabasePath = "financialtracker/"
	// CategorySchema is the name of the collection categories are stored in.
	CategorySchema = "category"
	// NameField is the category's unique name.
	NameField = "name"
	// CreatedAtField is the unix second the category was created at.
	CreatedAtField = "createdAt"
	// TransactionsField is the nested collection a category owns its
	// transactions in.
	TransactionsField = "transactions"
	// ReferenceField is a transaction's unique key inside its category.
	ReferenceField = "reference"
	// AmountField is a transaction's value in the smallest currency unit.
	AmountField = "amount"
	// KindField is a transaction's direction, api.Spend or api.Received.
	KindField = "kind"
	// OccurredAtField is the unix second a transaction was recorded at.
	OccurredAtField = "occurredAt"
	// ReferenceSeparator splits a transaction reference into its sequence
	// number and its description.
	ReferenceSeparator = "|"
)

// Props describes the tracker's database: one "category" collection whose
// records own a nested "transactions" collection. Keep offers unique string
// keys and integers, so a transaction's description — which is not unique —
// travels inside its unique reference; see Reference and Description.
func Props() keepdeps.Props {
	return keepdeps.Props{
		Path: DatabasePath,
		Schemas: []keepdeps.Schema{
			{
				Name: CategorySchema,
				Itens: []keepdeps.Item{
					{Name: NameField, Type: keepdeps.Key, Required: true},
					{Name: CreatedAtField, Type: keepdeps.Int, Required: true},
					{
						Name: TransactionsField,
						Type: keepdeps.Database,
						Itens: []keepdeps.Item{
							{Name: ReferenceField, Type: keepdeps.Key, Required: true},
							{Name: AmountField, Type: keepdeps.Int, Required: true},
							{Name: KindField, Type: keepdeps.Int, Required: true},
							{Name: OccurredAtField, Type: keepdeps.Int, Required: true},
						},
					},
				},
			},
		},
	}
}

// Categories returns the category collection of the injected database. ok is
// false when the injected library hands back no such schema.
func Categories(d deps.Deps) (keepdeps.SchemaInstance, bool) {
	return d.KeepLib.NewDatabase(Props()).GetSchema(CategorySchema)
}

// FindCategory returns the stored category record carrying the given name.
// ok is false when no category does.
func FindCategory(d deps.Deps, name string) (keepdeps.SchemaItem, bool) {
	categories, ok := Categories(d)
	if !ok {
		return keepdeps.SchemaItem{}, false
	}
	return categories.FindByKey(NameField, name)
}

// Text reads a Key field of a record, returning "" when the field is absent
// or does not hold a string.
func Text(record keepdeps.SchemaItem, field string) string {
	value, err := record.Get(field)
	if err != nil {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}

// Number reads an Int field of a record, returning 0 when the field is
// absent or does not hold an integer.
func Number(record keepdeps.SchemaItem, field string) int64 {
	value, err := record.Get(field)
	if err != nil {
		return 0
	}
	number, ok := value.(int64)
	if !ok {
		return 0
	}
	return number
}

// Reference composes a transaction's unique key from a sequence number and
// its description, so two transactions can carry the same description while
// the stored key stays unique.
func Reference(sequence int64, description string) string {
	return strconv.FormatInt(sequence, 10) + ReferenceSeparator + description
}

// Description reads the description back out of a transaction reference
// composed by Reference.
func Description(reference string) string {
	_, description, found := strings.Cut(reference, ReferenceSeparator)
	if !found {
		return reference
	}
	return description
}

// Money renders an amount held in the smallest currency unit as a decimal
// string with two places, keeping the sign of a negative amount.
func Money(amount int64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	cents := amount % 100
	padding := ""
	if cents < 10 {
		padding = "0"
	}
	return sign + strconv.FormatInt(amount/100, 10) + "." + padding + strconv.FormatInt(cents, 10)
}
