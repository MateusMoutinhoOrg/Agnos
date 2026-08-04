package cli

// The command-line interface of the tracker, written entirely inside the
// sandbox. It reads the command line through the injected Verb argv parser
// (deps.Deps.VerbLib) and writes every line through the injected formatted
// writer (deps.Deps.Printf), so the whole program stays free of OS-bound and
// third-party imports — the process only hands it an argument vector and
// exits with the code it returns.
//
// Like sandbox/internal/store, this package is neither an object package nor
// the entry point: it declares no types and no factories, and is called by
// SandboxmainFactory in sandbox/internal/lib.

import (
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/store"
)

// Version is the interface version reported by the `version` command and the
// `--version` flag.
const Version = "v0.1.0"

// Flag spellings the interface understands, in the shape Verb's IsPresent
// takes: every spelling of one flag in a single slice.
var (
	// HelpFlags asks for the usage screen instead of running a command.
	HelpFlags = []string{"-h", "--help"}
	// VersionFlags asks for the interface version instead of running a
	// command.
	VersionFlags = []string{"-v", "--version"}
	// QuietFlags suppress the confirmation lines a mutating command prints,
	// leaving only listings and errors.
	QuietFlags = []string{"-q", "--quiet"}
)

// Usage is the help screen, printed for `help`, for `--help`, and whenever a
// command line cannot be understood.
const Usage = `agnos — a financial tracker on the command line

Usage:
  agnos <command> [arguments] [flags]

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
  agnos category add groceries
  agnos spend groceries "weekly shopping" 84.50
  agnos received salary "august paycheck" 2500.00
  agnos balance groceries
`

// Run is the body of api.Lib.Sandboxmain: it dispatches one command line and
// returns the process exit code. args is only read to detect an empty
// command line; the arguments themselves are drained through the injected
// Verb parser, which the adapter wired over the same vector.
func Run(l *api.Lib, args []string) int {
	if len(args) == 0 {
		l.Deps.Printf("%s", Usage)
		return api.ExitUsage
	}

	verb := l.Deps.VerbLib
	if verb.IsPresent(HelpFlags) {
		l.Deps.Printf("%s", Usage)
		return api.ExitOk
	}
	if verb.IsPresent(VersionFlags) {
		return version(l)
	}
	// Read the flag before the positional arguments: Verb marks a matched
	// flag used, so draining what is left hands back only the command words.
	quiet := verb.IsPresent(QuietFlags)

	command, err := verb.GetNextStringArg()
	if err != nil {
		return usageError(l, "no command given")
	}

	switch command {
	case "help":
		l.Deps.Printf("%s", Usage)
		return api.ExitOk
	case "version":
		return version(l)
	case "category":
		return category(l, quiet)
	case "spend":
		return record(l, api.Spend, quiet)
	case "received":
		return record(l, api.Received, quiet)
	case "transactions":
		return transactions(l)
	case "balance":
		return balance(l)
	}
	return usageError(l, "unknown command "+quote(command))
}

// version prints the interface version. Both the `version` command and the
// --version flag land here.
func version(l *api.Lib) int {
	l.Deps.Printf("agnos %s\n", Version)
	return api.ExitOk
}

// category runs the `category` command group — add, list, remove — over the
// arguments still unread in the injected parser.
func category(l *api.Lib, quiet bool) int {
	verb := l.Deps.VerbLib
	action, err := verb.GetNextStringArg()
	if err != nil {
		return usageError(l, "category needs an action: add, list or remove")
	}

	switch action {
	case "list":
		categories := l.ListCategories()
		if len(categories) == 0 {
			l.Deps.Printf("no categories yet — create one with: agnos category add <name>\n")
			return api.ExitOk
		}
		for _, stored := range categories {
			l.Deps.Printf("%s\n", stored.String())
		}
		return api.ExitOk

	case "add":
		name, err := verb.GetNextStringArg()
		if err != nil {
			return usageError(l, "category add needs a name")
		}
		created, ok := l.AddCategory(name)
		if !ok {
			return failure(l, "could not create the category "+quote(name))
		}
		if !quiet {
			l.Deps.Printf("category %s\n", created.String())
		}
		return api.ExitOk

	case "remove":
		name, err := verb.GetNextStringArg()
		if err != nil {
			return usageError(l, "category remove needs a name")
		}
		stored, found := l.GetCategory(name)
		if !found {
			return failure(l, "no category named "+quote(name))
		}
		if !stored.Remove() {
			return failure(l, "could not remove the category "+quote(name))
		}
		if !quiet {
			l.Deps.Printf("removed category %s\n", name)
		}
		return api.ExitOk
	}
	return usageError(l, "unknown category action "+quote(action))
}

// record runs the `spend` and `received` commands, which differ only by the
// kind of transaction they write.
func record(l *api.Lib, kind int, quiet bool) int {
	verb := l.Deps.VerbLib
	name, nameErr := verb.GetNextStringArg()
	description, descriptionErr := verb.GetNextStringArg()
	amountText, amountErr := verb.GetNextStringArg()
	if nameErr != nil || descriptionErr != nil || amountErr != nil {
		return usageError(l, kindName(kind)+" needs a category, a description and an amount")
	}

	amount, ok := ParseAmount(amountText)
	if !ok {
		return usageError(l, "invalid amount "+quote(amountText)+" — expected a positive decimal like 84.50")
	}

	written := api.Transaction{}
	if kind == api.Spend {
		written, ok = l.AddSpend(name, description, amount)
	} else {
		written, ok = l.AddReceived(name, description, amount)
	}
	if !ok {
		return failure(l, "could not record the transaction under "+quote(name)+" — is the category created?")
	}
	if !quiet {
		l.Deps.Printf("%s\n", written.String())
	}
	return api.ExitOk
}

// transactions runs the `transactions` command, listing every stored
// transaction or, when a category name follows, only that category's.
func transactions(l *api.Lib) int {
	verb := l.Deps.VerbLib

	listed := []api.Transaction{}
	name, err := verb.GetNextStringArg()
	if err != nil {
		listed = l.ListTransactions()
	} else {
		stored, found := l.GetCategory(name)
		if !found {
			return failure(l, "no category named "+quote(name))
		}
		listed = stored.ListTransactions()
	}

	if len(listed) == 0 {
		l.Deps.Printf("no transactions yet — record one with: agnos spend <category> <description> <amount>\n")
		return api.ExitOk
	}
	for _, written := range listed {
		l.Deps.Printf("%s\n", written.String())
	}
	return api.ExitOk
}

// balance runs the `balance` command, printing the total across every
// category or, when a category name follows, that category's own balance.
func balance(l *api.Lib) int {
	verb := l.Deps.VerbLib

	name, err := verb.GetNextStringArg()
	if err != nil {
		l.Deps.Printf("%s\n", store.Money(l.Balance()))
		return api.ExitOk
	}

	stored, found := l.GetCategory(name)
	if !found {
		return failure(l, "no category named "+quote(name))
	}
	l.Deps.Printf("%s\n", store.Money(stored.Balance()))
	return api.ExitOk
}

// ParseAmount reads a decimal amount written the way a person types money —
// "84.50", "84.5", "84" — into the smallest currency unit the library
// stores. It reports false for anything that is not a positive decimal with
// at most two places, so the caller can answer with a usage error rather
// than silently record a wrong figure.
func ParseAmount(text string) (int64, bool) {
	units, cents, hasCents := strings.Cut(strings.TrimSpace(text), ".")
	if !digits(units) {
		return 0, false
	}
	whole, err := strconv.ParseInt(units, 10, 64)
	if err != nil {
		return 0, false
	}

	fraction := int64(0)
	if hasCents {
		if len(cents) == 1 {
			cents += "0"
		}
		if len(cents) != 2 || !digits(cents) {
			return 0, false
		}
		fraction, err = strconv.ParseInt(cents, 10, 64)
		if err != nil {
			return 0, false
		}
	}

	amount := whole*100 + fraction
	if amount <= 0 || whole > (1<<62)/100 {
		return 0, false
	}
	return amount, true
}

// digits reports whether text is a non-empty run of decimal digits, which is
// what keeps ParseAmount from accepting the signs and exponents strconv
// would otherwise take.
func digits(text string) bool {
	if text == "" {
		return false
	}
	for _, character := range text {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

// kindName renders a transaction kind as the command that records it.
func kindName(kind int) string {
	if kind == api.Spend {
		return "spend"
	}
	return "received"
}

// quote wraps a value read off the command line in quotes, so an empty or
// space-carrying argument is still visible in an error line.
func quote(value string) string {
	return `"` + value + `"`
}

// usageError reports a command line that could not be understood, printing
// the reason followed by the usage screen.
func usageError(l *api.Lib, reason string) int {
	l.Deps.Printf("agnos: %s\n\n", reason)
	l.Deps.Printf("%s", Usage)
	return api.ExitUsage
}

// failure reports a well-formed command that could not be carried out,
// printing the reason without the usage screen — the command line was fine,
// the records were not.
func failure(l *api.Lib, reason string) int {
	l.Deps.Printf("agnos: %s\n", reason)
	return api.ExitFailure
}
