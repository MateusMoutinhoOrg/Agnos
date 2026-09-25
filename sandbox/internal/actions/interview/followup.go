package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// A command of agnos declares one thing; the thing itself is the several
// commands that fill it in. `add-route` writes a route that answers nothing
// and carries no body, and the commands that give it one are six rows down a
// menu the session had already gone back to — so the person who has just
// declared a route has to know that `add-body-field` exists before they can
// look for it.
//
// This file is the answer: what a command leaves half-declared, and what
// finishes it. It is the fourth table of the agnos vocabulary this package
// spells, beside SuggestFor, RuledOut and the gate of state.go, and like them
// it names commands rather than reading them off a declaration — a route is
// finished by its fields because of what a route is, which nothing in
// entries.yaml says.

// followUp is one command worth running next: the command itself, what it does
// in the words of someone who has just watched the one before it finish, and
// the answer it inherits — the field of the finished command that supplies the
// field of this one, so the route just declared is not typed again.
type followUp struct {
	Verb string
	Msg  string
	From string
	To   string
}

// followUps is what comes after each command that leaves something
// half-declared. The order is the order the thing is usually filled in, and
// the first row is what most sessions want next — a route is declared to carry
// a body far more often than it is declared to read a header.
//
// A command that finishes what it started has no entry, and a session that
// runs one goes straight back to the menu.
var followUps = map[string][]followUp{
	"add-route": {
		{"import-body", "Read its body from an example payload", "name", "route"},
		{"add-body-field", "Declare one property of its body", "name", "route"},
		{"set-body", "Say what kind of body it takes", "name", "route"},
		{"add-parameter", "Declare a value it reads from the query or a header", "name", "route"},
		{"add-path", "Read one more slice of its path", "name", "route"},
		{"show-route", "Look at what it declares so far", "name", "route"},
	},
	"add-page": {
		{"add-parameter", "Declare a value it reads from the query or a header", "name", "route"},
		{"add-path", "Read one more slice of its path", "name", "route"},
		{"show-route", "Look at what it declares so far", "name", "route"},
	},
	"import-body": {
		{"show-route", "Look at what the route declares now", "route", "route"},
		{"set-body-field", "Change one of the properties it read", "route", "route"},
		{"add-body-field", "Declare one the example did not carry", "route", "route"},
	},
	"add-body-field": {
		{"add-body-field", "Declare one more property", "route", "route"},
		{"show-route", "Look at what the route declares now", "route", "route"},
	},
	"set-body-field": {
		{"show-route", "Look at what the route declares now", "route", "route"},
		{"set-body-field", "Change one more property", "route", "route"},
	},
	"show-route": {
		{"add-body-field", "Declare one more property of its body", "route", "route"},
		{"set-body-field", "Change one of its body properties", "route", "route"},
		{"add-parameter", "Declare a value it reads from the query or a header", "route", "route"},
		{"set-route", "Change what the route itself says", "route", "route"},
	},
	"set-body": {
		{"add-body-field", "Declare one property of its body", "route", "route"},
		{"show-route", "Look at what the route declares now", "route", "route"},
	},
	"add-parameter": {
		{"add-parameter", "Declare one more value it reads", "route", "route"},
		{"show-route", "Look at what the route declares now", "route", "route"},
	},
	"add-path": {
		{"add-path", "Read one more slice of its path", "route", "route"},
		{"show-route", "Look at what the route declares now", "route", "route"},
	},
	"set-parameter": {
		{"show-route", "Look at what the route declares now", "route", "route"},
		{"set-parameter", "Change one more value it reads", "route", "route"},
	},
	"set-path": {
		{"show-route", "Look at what the route declares now", "route", "route"},
		{"set-path", "Change one more slice of its path", "route", "route"},
	},
	"set-route": {
		{"show-route", "Look at what the route declares now", "route", "route"},
	},
	"add-database": {
		{"add-table", "Declare a table it holds", "name", "database"},
		{"show-database", "Look at what it declares so far", "name", "database"},
	},
	"add-table": {
		{"add-table-field", "Declare one field of it", "name", "table"},
		{"add-table", "Declare one more table", "database", "database"},
		{"show-database", "Look at what the database declares now", "database", "database"},
	},
	"add-table-field": {
		{"add-table-field", "Declare one more field", "table", "table"},
		{"show-database", "Look at what the database declares now", "database", "database"},
	},
	"set-table-field": {
		{"show-database", "Look at what the database declares now", "database", "database"},
		{"set-table-field", "Change one more field", "table", "table"},
	},
	"show-database": {
		{"add-table", "Declare one more table", "database", "database"},
	},
	"add-command": {
		{"add-flag", "Give it a flag", "name", "command"},
		{"add-arg", "Give it an argument", "name", "command"},
		{"set-command", "Say more about it — a long description, an example", "name", "command"},
	},
	"add-flag": {
		{"add-flag", "Give it one more flag", "command", "command"},
		{"add-arg", "Give it an argument", "command", "command"},
	},
	"add-arg": {
		{"add-arg", "Give it one more argument", "command", "command"},
		{"add-flag", "Give it a flag", "command", "command"},
	},
}

// carriedNames are the commands whose --name is written down under a spelling
// of agnos's choosing. What a follow-up inherits is the name on disk, never
// the one that was typed: `add-route "Create User"` declares create-user, and
// a follow-up carrying "Create User" would name a route that does not exist.
var carriedNames = map[string]bool{
	"add-route":    true,
	"add-page":     true,
	"add-command":  true,
	"add-database": true,
	"add-table":    true,
}

// chooseFollowUp offers what comes after the command that has just finished,
// and reports false when the person chose to go back to the menu instead. The
// answers it carries are already on the confirm screen of the command it
// picks, so nothing about the route or the command in hand is asked twice.
func chooseFollowUp(sandbox *api.Sandbox, io *smartio.SmartIO, finished api.Command, values map[string][]any) (api.Command, map[string][]any, bool, error) {
	offers := offeredFollowUps(sandbox, io, finished, values)
	if len(offers) == 0 {
		return api.Command{}, nil, false, nil
	}

	rows := []interviewer.AlternativeOption{}
	for _, offer := range offers {
		rows = append(rows, interviewer.AlternativeOption{
			Id:  offer.Verb,
			Msg: sandbox.Deps.Std.Sprintf("%s  %s(%s)%s", offer.Msg, dim, offer.Verb, reset),
		})
	}
	rows = append(rows, interviewer.AlternativeOption{Id: backOptionId, Msg: "· nothing else — back to the menu"})

	chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion(followUpQuestion(sandbox, finished, offers[0], values), rows)
	if err != nil {
		// Going back from here is the same answer as the last row: the
		// command has already run, so there is nothing to undo.
		if sandbox.Deps.Interviewer.Back(err) {
			return api.Command{}, nil, false, nil
		}
		return api.Command{}, nil, false, err
	}
	if chosen == backOptionId {
		return api.Command{}, nil, false, nil
	}

	for _, offer := range offers {
		if offer.Verb != chosen {
			continue
		}
		command, found := commandByVerb(sandbox, offer.Verb)
		if !found {
			return api.Command{}, nil, false, nil
		}
		return command, carried(sandbox, finished, command, offer, values), true, nil
	}

	return api.Command{}, nil, false, nil
}

// offeredFollowUps is the table filtered by what this project can run and by
// what the finished command actually answered: a follow-up inheriting an
// answer that was skipped has nothing to carry, and offering it would ask for
// the route by hand under a heading that promised not to.
func offeredFollowUps(sandbox *api.Sandbox, io *smartio.SmartIO, finished api.Command, values map[string][]any) []followUp {
	state := readState(sandbox, io)

	offers := []followUp{}
	for _, offer := range followUps[verbOf(finished)] {
		if answeredText(sandbox, values, offer.From) == "" {
			continue
		}
		command, found := commandByVerb(sandbox, offer.Verb)
		if !found || !offerable(state, command) {
			continue
		}
		offers = append(offers, offer)
	}

	return offers
}

// followUpQuestion names the thing the follow-ups are about, so the menu reads
// as being about the route just declared rather than about commands in
// general. Every offer of one list inherits the same answer, so the first of
// them names it.
func followUpQuestion(sandbox *api.Sandbox, finished api.Command, first followUp, values map[string][]any) string {
	subject := carriedText(sandbox, finished, first, values)
	if subject == "" {
		return "What next?"
	}
	return sandbox.Deps.Std.Sprintf("What next for %s?", subject)
}

// carried is what the follow-up starts with: the answer it inherits, bound
// under the field that follow-up declares it as, plus every scope field the
// finished command answered that the next one declares too. The second half is
// what makes a two-deep unit work — a field of a table is named by its
// database and its table both, and only one of the two can be the inherited
// answer.
//
// A scope field the next command does not declare is left out: what is bound
// here is read back as that command's own answers, so it may only hold ids
// that command has.
func carried(sandbox *api.Sandbox, finished api.Command, next api.Command, offer followUp, values map[string][]any) map[string][]any {
	inherited := map[string][]any{}

	for _, id := range scopeFields {
		text := answeredText(sandbox, values, id)
		if text == "" || !declaresField(next, id) {
			continue
		}
		inherited[id] = []any{text}
	}

	if text := carriedText(sandbox, finished, offer, values); text != "" {
		inherited[offer.To] = []any{text}
	}

	return inherited
}

// declaresField reports whether one command declares a field under this id.
func declaresField(command api.Command, id string) bool {
	for _, field := range FieldsOf(command) {
		if field.Id == id {
			return true
		}
	}
	return false
}

// carriedText is one inherited answer as the value the next command is given:
// the name on disk when the finished command wrote a spelling of its own, and
// what was answered otherwise.
func carriedText(sandbox *api.Sandbox, finished api.Command, offer followUp, values map[string][]any) string {
	text := answeredText(sandbox, values, offer.From)
	if text == "" {
		return ""
	}
	if offer.From == nameFieldId && carriedNames[verbOf(finished)] {
		return utils.CommandIdentifier(sandbox, text)
	}
	return text
}
