package commands

import (
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// Record runs the `spend` and `received` commands, which differ only by the
// kind of transaction they write.
func Record(l *api.Lib, kind int, quiet bool) int {
	verb := l.Deps.VerbLib
	name, nameErr := verb.GetNextStringArg()
	description, descriptionErr := verb.GetNextStringArg()
	amountText, amountErr := verb.GetNextStringArg()
	if nameErr != nil || descriptionErr != nil || amountErr != nil {
		return UsageError(l, config.RecordOperandsMissing, kindName(kind))
	}

	amount, ok := ParseAmount(amountText)
	if !ok {
		return UsageError(l, config.AmountInvalid, amountText)
	}

	written := api.Transaction{}
	if kind == api.Spend {
		written, ok = l.AddSpend(name, description, amount)
	} else {
		written, ok = l.AddReceived(name, description, amount)
	}
	if !ok {
		return Failure(l, config.TransactionNotRecorded, name)
	}
	if !quiet {
		l.Deps.Printf("%s\n", written.String())
	}
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
