package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The field ids the declaration commands — add-flag, add-arg, add-param,
// add-header, add-segment, add-body-field and set-body — constrain each other
// by. No other command of the surface declares one of them, which is what lets
// the table below be keyed by id alone.
const (
	typeFieldId              = "type"
	defaultFieldId           = "default"
	requiredFieldId          = "required"
	optionalFieldId          = "optional"
	arrayFieldId             = "array"
	minFieldId               = "min"
	maxFieldId               = "max"
	exclusiveMinFieldId      = "exclusive-min"
	exclusiveMaxFieldId      = "exclusive-max"
	formatFieldId            = "format"
	patternFieldId           = "pattern"
	minItemsFieldId          = "min-items"
	maxItemsFieldId          = "max-items"
	uniqueItemsFieldId       = "unique-items"
	additionalPropsFieldId   = "additional-properties"
	noAdditionalPropsFieldId = "no-additional-properties"
	targetFieldId            = "target"
)

// typeObject is the one declared type that is a json-schema kind rather than a
// value the cli converts, so add-body-field is the only command that sees it.
const typeObject = "object"

// The two declared field types of a table that the rules below read: a link is
// the only one that points at another table, and a nested collection is the
// only one an insert never writes.
const (
	typeLink          = "link"
	typeNestedRecords = "database"
)

// tableFieldVerbs are the commands declaring a field of a database table.
// --target and --required read differently there — a target belongs to a link
// alone, and a nested collection is never required — and declaring one and
// editing one are the same declaration, so both are here.
var tableFieldVerbs = map[string]bool{
	"add-table-field": true,
	"set-table-field": true,
}

// schemaVerbs are the commands declaring a json-schema property instead of an
// agnos field. Two rules read differently there: a required boolean property
// is legal — required means the key has to be present, not that a false is
// demanded — and min/max also bound the length of a string. Declaring one and
// editing one are the same declaration, so both are here.
var schemaVerbs = map[string]bool{
	"add-body-field": true,
	"set-body-field": true,
}

// editVerbs are the commands that rewrite a declaration instead of writing
// one. They are the one place an unanswered field is not a fact: everywhere
// else a skipped --type means the declared default and a no to --array means
// not a list, while here both mean "leave it as it is" — and the type left as
// it is may well be the one the keyword applies to.
var editVerbs = map[string]bool{
	"set-segment":    true,
	"set-header":     true,
	"set-param":      true,
	"set-body-field": true,
}

// RuledOut reports a field there is nothing left to ask about, because an
// answer already given decides it. Asking anyway offers a combination the
// command refuses — a flag both required and defaulted, a min on a string —
// and the person only finds out after the confirm screen.
//
// It is the second half of what the interview knows about agnos's own
// vocabulary, beside SuggestFor: that table says where a field's answers come
// from, this one which fields a previous answer removes. Both are keyed by
// field id, and a field with no entry here is always asked.
func RuledOut(sandbox *api.Sandbox, command api.Command, field Field, values map[string][]any) bool {
	return RuledOutReason(sandbox, command, field, values) != ""
}

// RuledOutReason is why a field was not asked, and "" for one that was. The
// rule and the sentence are one thing: a question that disappears without a
// word is the interview's own version of a hidden flag, so every rule below
// answers with what it took away and the confirm screen says so.
//
// Every one of them is a rule the action behind the command enforces with an
// error: utils.NewField and utils.NewRouteField for the fields,
// utils.RouteBodyPropertySchema for the schema keywords, set-body for its own
// pair.
func RuledOutReason(sandbox *api.Sandbox, command api.Command, field Field, values map[string][]any) string {
	kind := answeredText(sandbox, values, typeFieldId)
	schema := schemaVerbs[verbOf(command)]

	// An editor was told the type only if it was answered, and on every other
	// command a skipped --type binds the declared default — so an empty kind
	// means the type is whatever is already on disk, and a rule that reads one
	// has nothing to read. Rules that read no type are unaffected.
	typed := kind != ""
	listed := answeredYes(values, arrayFieldId) || editVerbs[verbOf(command)]

	table := tableFieldVerbs[verbOf(command)]

	switch field.Id {
	case targetFieldId:
		if table && typed && kind != typeLink {
			return "only a reference to another table points at one"
		}
	case requiredFieldId:
		if table {
			if typed && kind == typeNestedRecords {
				return "a collection nested under each record is never written when one is inserted"
			}
			break
		}
		if answeredText(sandbox, values, defaultFieldId) != "" {
			return "what it falls back to already covers its absence"
		}
		if typed && !schema && kind == typeBoolean {
			return "a yes-or-no value is never demanded — not giving it already means no"
		}
	case defaultFieldId:
		if answeredYes(values, requiredFieldId) {
			return "it has to be given, so there is nothing to fall back to"
		}
	case optionalFieldId:
		if answeredYes(values, requiredFieldId) {
			return "it has to be given, so there is nothing to fall back to"
		}
	case minFieldId, maxFieldId:
		if typed && schema && (kind == typeBoolean || kind == typeObject) {
			return "only a number or a piece of text carries a smallest and a largest"
		}
		if typed && !schema && !numeric(kind) {
			return "only a number carries a smallest and a largest"
		}
	case exclusiveMinFieldId, exclusiveMaxFieldId:
		if typed && !numeric(kind) {
			return "only a number carries a bound it may not touch"
		}
	case formatFieldId, patternFieldId:
		if typed && kind != typeString {
			return "only a piece of text carries a shape"
		}
	case minItemsFieldId, maxItemsFieldId, uniqueItemsFieldId:
		if !listed {
			return "only a list carries a length"
		}
	case additionalPropsFieldId:
		if typed && kind != typeObject {
			return "only an object holds keys of its own"
		}
	case noAdditionalPropsFieldId:
		if typed && kind != typeObject {
			return "only an object holds keys of its own"
		}
		if answeredYes(values, additionalPropsFieldId) {
			return "the question before this one already answered it"
		}
	}

	return ""
}

// NotAskedNotes is what the confirm screen says about the questions that never
// appeared: one line per reason, naming every field it took away. The screen
// promises the line is what the person could have typed themselves, and a flag
// they looked for and never saw asked about is the one thing that makes that
// line read as a shorter command than it is.
func NotAskedNotes(sandbox *api.Sandbox, command api.Command, values map[string][]any) []string {
	reasons := []string{}
	taken := map[string][]string{}

	for _, field := range FieldsOf(command) {
		reason := RuledOutReason(sandbox, command, field, values)
		if reason == "" {
			continue
		}
		if _, seen := taken[reason]; !seen {
			reasons = append(reasons, reason)
		}
		taken[reason] = append(taken[reason], label(field))
	}

	notes := []string{}
	for _, reason := range reasons {
		notes = append(notes, sandbox.Deps.Std.Sprintf("not asked: %s — %s",
			sandbox.Deps.Stringsdeps.Join(taken[reason], ", "), reason))
	}

	return notes
}

// PruneRuledOut drops every answer a change on the confirm screen has just
// made meaningless — the min typed while the type was int, after the type
// became string. It walks the fields in the order they are asked, so a rule
// reads the same answers it read during the session.
func PruneRuledOut(sandbox *api.Sandbox, command api.Command, values map[string][]any) {
	for _, field := range AskOrder(FieldsOf(command)) {
		if RuledOut(sandbox, command, field, values) {
			delete(values, field.Id)
		}
	}
}

// numeric reports the two declared types a bound applies to. Every command
// above declares "string" as its --type default, so an unanswered type is
// bound to it and never reads back empty here.
func numeric(kind string) bool {
	return kind == typeInt || kind == typeFloat
}

// answeredText is one answer already given, written the way a command line
// carries it. A field nobody answered — or skipped, which binds nothing —
// reads back as "".
func answeredText(sandbox *api.Sandbox, values map[string][]any, id string) string {
	bound := values[id]
	if len(bound) == 0 {
		return ""
	}
	return valueText(sandbox, bound[0])
}

// answeredYes reports a yes-or-no field answered yes.
func answeredYes(values map[string][]any, id string) bool {
	bound := values[id]
	if len(bound) == 0 {
		return false
	}
	truth, ok := bound[0].(bool)
	return ok && truth
}
