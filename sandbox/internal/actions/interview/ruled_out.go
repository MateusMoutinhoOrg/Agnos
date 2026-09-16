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
)

// typeObject is the one declared type that is a json-schema kind rather than a
// value the cli converts, so add-body-field is the only command that sees it.
const typeObject = "object"

// schemaVerb is the command declaring a json-schema property instead of an
// agnos field. Two rules read differently there: a required boolean property
// is legal — required means the key has to be present, not that a false is
// demanded — and min/max also bound the length of a string.
const schemaVerb = "add-body-field"

// RuledOut reports a field there is nothing left to ask about, because an
// answer already given decides it. Asking anyway offers a combination the
// command refuses — a flag both required and defaulted, a min on a string —
// and the person only finds out after the confirm screen.
//
// It is the second half of what the interview knows about agnos's own
// vocabulary, beside SuggestFor: that table says where a field's answers come
// from, this one which fields a previous answer removes. Both are keyed by
// field id, and a field with no entry here is always asked.
//
// Every rule is one the action behind the command enforces with an error:
// utils.NewField and utils.NewRouteField for the fields, add-body-field's
// propertySchema for the schema keywords, set-body for its own pair.
func RuledOut(sandbox *api.Sandbox, command api.Command, field Field, values map[string][]any) bool {
	kind := answeredText(sandbox, values, typeFieldId)
	schema := verbOf(command) == schemaVerb

	switch field.Id {
	case requiredFieldId:
		return answeredText(sandbox, values, defaultFieldId) != "" || (!schema && kind == typeBoolean)
	case defaultFieldId:
		return answeredYes(values, requiredFieldId)
	case optionalFieldId:
		return answeredYes(values, requiredFieldId)
	case minFieldId, maxFieldId:
		if schema {
			return kind == typeBoolean || kind == typeObject
		}
		return !numeric(kind)
	case exclusiveMinFieldId, exclusiveMaxFieldId:
		return !numeric(kind)
	case formatFieldId, patternFieldId:
		return kind != typeString
	case minItemsFieldId, maxItemsFieldId, uniqueItemsFieldId:
		return !answeredYes(values, arrayFieldId)
	case additionalPropsFieldId:
		return kind != typeObject
	case noAdditionalPropsFieldId:
		return kind != typeObject || answeredYes(values, additionalPropsFieldId)
	}

	return false
}

// PruneRuledOut drops every answer a change on the confirm screen has just
// made meaningless — the min typed while the type was int, after the type
// became string. It walks the fields in the order they are asked, so a rule
// reads the same answers it read during the session.
func PruneRuledOut(sandbox *api.Sandbox, command api.Command, values map[string][]any) {
	for _, field := range FieldsOf(command) {
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
