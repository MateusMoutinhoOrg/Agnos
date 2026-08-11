package cli

// The command-line interface of the tracker, written entirely inside the
// sandbox. It reads the command line through the injected Verb argv parser
// (deps.Deps.VerbLib), takes every word it displays from the injected
// embedded assets (deps.Deps.EmbedDeps), and writes every line through the
// injected formatted writer (deps.Deps.Printf), so the whole program stays
// free of OS-bound and third-party imports — the process only hands it an
// argument vector and exits with the code it returns.
//
// No display text is written here: the usage screen, the version, and each
// message below live in files under /assets/, and this package addresses them
// by path. Changing what the interface says is editing an asset, not editing
// Go.
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

// Paths of the assets the interface reads, relative to the root of the asset
// tree the adapter serves through deps.Deps.EmbedDeps.
const (
	// VersionAsset holds the interface version reported by the `version`
	// command and the `--version` flag — the single place a release bump is
	// written.
	VersionAsset = "version.txt"
	// UsageAsset holds the help screen, printed for `help`, for `--help`,
	// and whenever a command line cannot be understood.
	UsageAsset = "usages.txt"
	// MessagesDir holds one file per line the interface can print, each
	// named after what it reports and read through message.
	MessagesDir = "messages/"
)

// Names of the messages under MessagesDir, one per line the interface can
// print. A message is a Printf format: the ones naming a value carry the
// verb — and the quotes — the value is rendered with.
const (
	// ErrorPrefixMessage opens every line reporting a refused command.
	ErrorPrefixMessage = "error-prefix"
	// VersionMessage renders the interface version as one line.
	VersionMessage = "version"
	// CategoryAddedMessage confirms a created category.
	CategoryAddedMessage = "category-added"
	// CategoryRemovedMessage confirms a deleted category.
	CategoryRemovedMessage = "category-removed"
	// NoCategoriesMessage answers a listing with nothing stored yet.
	NoCategoriesMessage = "no-categories"
	// NoTransactionsMessage answers a transaction listing with nothing
	// recorded yet.
	NoTransactionsMessage = "no-transactions"

	// NoCommandMessage reports a command line carrying no command at all.
	NoCommandMessage = "no-command"
	// UnknownCommandMessage reports a command the interface does not have.
	UnknownCommandMessage = "unknown-command"
	// CategoryActionMissingMessage reports `category` with no action after
	// it.
	CategoryActionMissingMessage = "category-action-missing"
	// CategoryActionUnknownMessage reports a `category` action that does not
	// exist.
	CategoryActionUnknownMessage = "category-action-unknown"
	// CategoryAddNameMissingMessage reports `category add` with no name.
	CategoryAddNameMissingMessage = "category-add-name-missing"
	// CategoryRemoveNameMissingMessage reports `category remove` with no
	// name.
	CategoryRemoveNameMissingMessage = "category-remove-name-missing"
	// RecordOperandsMissingMessage reports `spend` or `received` missing one
	// of its three operands.
	RecordOperandsMissingMessage = "record-operands-missing"
	// AmountInvalidMessage reports an amount ParseAmount refused.
	AmountInvalidMessage = "amount-invalid"

	// CategoryNotCreatedMessage reports a category that could not be
	// written.
	CategoryNotCreatedMessage = "category-not-created"
	// CategoryNotFoundMessage reports a named category nothing is stored
	// under.
	CategoryNotFoundMessage = "category-not-found"
	// CategoryNotRemovedMessage reports a category that could not be
	// deleted.
	CategoryNotRemovedMessage = "category-not-removed"
	// TransactionNotRecordedMessage reports a transaction that could not be
	// written.
	TransactionNotRecordedMessage = "transaction-not-recorded"
)

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

// asset reads one embedded asset through the injected embed library and
// returns it verbatim, newlines and all. An asset that cannot be read is a
// packaging mistake rather than a user mistake, so it reports itself by path
// instead of printing nothing.
func asset(l *api.Lib, path string) string {
	content, err := l.Deps.EmbedDeps.ReadFile(path)
	if err != nil {
		return "agnos-cli: missing asset " + path + "\n"
	}
	return string(content)
}

// message reads one line the interface prints, named after what it reports.
// The trailing newline every asset file ends with is trimmed, so the call
// site decides the layout — the same message can close a line, open a
// paragraph, or be followed by the usage screen.
func message(l *api.Lib, name string) string {
	return strings.TrimRight(asset(l, MessagesDir+name+".txt"), "\r\n")
}

// Run is the body of api.Lib.Sandboxmain: it dispatches one command line and
// returns the process exit code. args is only read to detect an empty
// command line; the arguments themselves are drained through the injected
// Verb parser, which the adapter wired over the same vector.
func Run(l *api.Lib, args []string) int {
	if len(args) == 0 {
		l.Deps.Printf("%s", asset(l, UsageAsset))
		return api.ExitUsage
	}

	verb := l.Deps.VerbLib
	if verb.IsPresent(HelpFlags) {
		l.Deps.Printf("%s", asset(l, UsageAsset))
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
		return usageError(l, NoCommandMessage)
	}

	switch command {
	case "help":
		l.Deps.Printf("%s", asset(l, UsageAsset))
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
	return usageError(l, UnknownCommandMessage, command)
}

// version prints the interface version, read from its own asset so a release
// bump is a one-line edit outside the code. Both the `version` command and
// the --version flag land here.
func version(l *api.Lib) int {
	l.Deps.Printf(message(l, VersionMessage)+"\n", Version(l))
	return api.ExitOk
}

// Version returns the interface version the `version` command reports,
// trimmed of the newline its asset file ends with.
func Version(l *api.Lib) string {
	return strings.TrimSpace(asset(l, VersionAsset))
}

// category runs the `category` command group — add, list, remove — over the
// arguments still unread in the injected parser.
func category(l *api.Lib, quiet bool) int {
	verb := l.Deps.VerbLib
	action, err := verb.GetNextStringArg()
	if err != nil {
		return usageError(l, CategoryActionMissingMessage)
	}

	switch action {
	case "list":
		categories := l.ListCategories()
		if len(categories) == 0 {
			l.Deps.Printf("%s\n", message(l, NoCategoriesMessage))
			return api.ExitOk
		}
		for _, stored := range categories {
			l.Deps.Printf("%s\n", stored.String())
		}
		return api.ExitOk

	case "add":
		name, err := verb.GetNextStringArg()
		if err != nil {
			return usageError(l, CategoryAddNameMissingMessage)
		}
		created, ok := l.AddCategory(name)
		if !ok {
			return failure(l, CategoryNotCreatedMessage, name)
		}
		if !quiet {
			l.Deps.Printf(message(l, CategoryAddedMessage)+"\n", created.String())
		}
		return api.ExitOk

	case "remove":
		name, err := verb.GetNextStringArg()
		if err != nil {
			return usageError(l, CategoryRemoveNameMissingMessage)
		}
		stored, found := l.GetCategory(name)
		if !found {
			return failure(l, CategoryNotFoundMessage, name)
		}
		if !stored.Remove() {
			return failure(l, CategoryNotRemovedMessage, name)
		}
		if !quiet {
			l.Deps.Printf(message(l, CategoryRemovedMessage)+"\n", name)
		}
		return api.ExitOk
	}
	return usageError(l, CategoryActionUnknownMessage, action)
}

// record runs the `spend` and `received` commands, which differ only by the
// kind of transaction they write.
func record(l *api.Lib, kind int, quiet bool) int {
	verb := l.Deps.VerbLib
	name, nameErr := verb.GetNextStringArg()
	description, descriptionErr := verb.GetNextStringArg()
	amountText, amountErr := verb.GetNextStringArg()
	if nameErr != nil || descriptionErr != nil || amountErr != nil {
		return usageError(l, RecordOperandsMissingMessage, kindName(kind))
	}

	amount, ok := ParseAmount(amountText)
	if !ok {
		return usageError(l, AmountInvalidMessage, amountText)
	}

	written := api.Transaction{}
	if kind == api.Spend {
		written, ok = l.AddSpend(name, description, amount)
	} else {
		written, ok = l.AddReceived(name, description, amount)
	}
	if !ok {
		return failure(l, TransactionNotRecordedMessage, name)
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
			return failure(l, CategoryNotFoundMessage, name)
		}
		listed = stored.ListTransactions()
	}

	if len(listed) == 0 {
		l.Deps.Printf("%s\n", message(l, NoTransactionsMessage))
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
		return failure(l, CategoryNotFoundMessage, name)
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

// usageError reports a command line that could not be understood, printing
// the named message — with whatever values it names rendered into it — and
// then the usage screen. The message asset carries the quotes a value is
// shown in, so an empty or space-carrying argument is still visible.
func usageError(l *api.Lib, name string, a ...any) int {
	l.Deps.Printf("%s ", message(l, ErrorPrefixMessage))
	l.Deps.Printf(message(l, name)+"\n\n", a...)
	l.Deps.Printf("%s", asset(l, UsageAsset))
	return api.ExitUsage
}

// failure reports a well-formed command that could not be carried out,
// printing the named message without the usage screen — the command line was
// fine, the records were not.
func failure(l *api.Lib, name string, a ...any) int {
	l.Deps.Printf("%s ", message(l, ErrorPrefixMessage))
	l.Deps.Printf(message(l, name)+"\n", a...)
	return api.ExitFailure
}
