package config

const (
	Usages = `agnos-cli — a financial tracker on the command line

Usage:
  agnos-cli <command> [arguments] [flags]

Commands:
  category add <name>                           create a category
  category list                                 list every category with its balance
  category remove <name>                        delete a category and its transactions
  spend <category> <description> <amount>       record money leaving the budget
  received <category> <description> <amount>    record money entering the budget
  transactions [category]                       list transactions, all or of one category
  balance [category]                            print the balance, total or of one category
  help                                          print this screen
  version                                       print the interface version

Flags:
  -h, --help                                    print this screen and exit
  -v, --version                                 print the interface version and exit
  -q, --quiet                                   print only listings and errors

Amounts are decimal, with at most two places: 84.50, 84.5 and 84 are all
accepted. They are always positive — the command chooses the direction.

Examples:
  agnos-cli category add groceries
  agnos-cli spend groceries "weekly shopping" 84.50
  agnos-cli received salary "august paycheck" 2500.00
  agnos-cli balance groceries
`
	CategoryActionUnknown = `unknown category action "%s"`
	CategoryAdded = `category %s`
	CategoryNotFound = `no category named "%s"`
	UnknownCommand = `unknown command "%s"`
	CategoryRemoveNameMissing = `category remove needs a name`
	TransactionNotRecorded = `could not record the transaction under "%s" — is the category created?`
	NoCategories = `no categories yet — create one with: agnos-cli category add <name>`
	CategoryRemoved = `removed category %s`
	NoTransactions = `no transactions yet — record one with: agnos-cli spend <category> <description> <amount>`
	CategoryAddNameMissing = `category add needs a name`
	CategoryNotCreated = `could not create the category "%s"`
	RecordOperandsMissing = `%s needs a category, a description and an amount`
	AmountInvalid = `invalid amount "%s" — expected a positive decimal like 84.50`
	CategoryActionMissing = `category needs an action: add, list or remove`
	ErrorPrefix = `agnos-cli:`
	CategoryNotRemoved = `could not remove the category "%s"`
	VersionMessage = `agnos-cli %s`
	NoCommand = `no command given`
)
