package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	interviewer "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/interviewer"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// The declared types a flag or an arg may carry, spelled as entries.yaml
// spells them — the same four the cli dispatch converts.
const (
	typeBoolean = "boolean"
	typeInt     = "int"
	typeFloat   = "float"
)

// The two fields every agnos command declares and the interview answers for
// the person: path is the project the session was opened on, and quiet stays
// off so the progress of a build is visible while it runs.
const (
	pathFieldId  = "path"
	quietFieldId = "quiet"
)

// scaffoldVerb is the one command whose --path is not the project being worked
// on but the directory a new one is written into. Binding the session's path
// to it would aim every scaffold at the project already open, so it is asked
// for like any other value.
const scaffoldVerb = "start"

// The two rows a menu may carry beside the real options. They are spelled with
// a leading NUL so that no value a project could declare collides with them.
const (
	skipOptionId  = "\x00skip"
	otherOptionId = "\x00other"
)

// Field is one flag or one arg, read the same way. A command declares the two
// in separate slices with almost the same keys, and every question below cares
// about the keys they share — so they are normalized into one shape once,
// rather than asked about twice.
type Field struct {
	Id          string
	Type        string
	Required    bool
	Array       bool
	Description string
	Default     string
	HasDefault  bool
	Min         float64
	HasMin      bool
	Max         float64
	HasMax      bool
	Identifiers []string
	IsFlag      bool
}

// FieldsOf is every value a command takes, args first and flags after — the
// order the person reading a help screen meets them in, and the order that
// puts the subject of the command before the options on it.
func FieldsOf(command api.Command) []Field {
	fields := []Field{}

	for _, arg := range command.Args {
		fields = append(fields, Field{
			Id: arg.Id, Type: arg.Type, Required: arg.Required, Array: arg.Array,
			Description: arg.Description, Default: arg.Default, HasDefault: arg.HasDefault,
			Min: arg.Min, HasMin: arg.HasMin, Max: arg.Max, HasMax: arg.HasMax,
			IsFlag: false,
		})
	}

	for _, flag := range command.Flags {
		fields = append(fields, Field{
			Id: flag.Id, Type: flag.Type, Required: flag.Required, Array: flag.Array,
			Description: flag.Description, Default: flag.Default, HasDefault: flag.HasDefault,
			Min: flag.Min, HasMin: flag.HasMin, Max: flag.Max, HasMax: flag.HasMax,
			Identifiers: flag.Identifiers, IsFlag: true,
		})
	}

	return fields
}

// AskValues asks one question per declared field and returns what to bind,
// keyed by field id. A field left unanswered carries its declared default
// instead — the dispatch does that when it reads a command line, and nothing
// else would, because the interview hands the handler its values directly.
func AskValues(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, session string) (map[string][]any, error) {
	values := map[string][]any{}

	for _, field := range FieldsOf(command) {
		bound, err := ResolveField(sandbox, io, command, field, session)
		if err != nil {
			return nil, err
		}
		if len(bound) > 0 {
			values[field.Id] = bound
		}
	}

	return values, nil
}

// ResolveField settles one field: the two the interview answers by itself, and
// everything else by asking. It is also what the confirm screen calls to
// change a single answer without walking the whole command again.
func ResolveField(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, field Field, session string) ([]any, error) {
	if AnsweredForYou(command, field) {
		if field.Id == pathFieldId {
			return []any{session}, nil
		}
		return []any{false}, nil
	}

	// The scaffold's --path is the one field whose declared default is not the
	// right offer: the session was opened on a folder, and a project created
	// from that session belongs in it.
	if field.Id == pathFieldId && verbOf(command) == scaffoldVerb {
		field.Default = session
		field.HasDefault = true
	}

	values, err := askField(sandbox, io, command, field)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 && field.HasDefault {
		return []any{DefaultValue(sandbox, field)}, nil
	}

	return values, nil
}

// AnsweredForYou reports the fields the session settles without asking: the
// project it was opened on, and the progress channel it keeps open. A scaffold
// is asked for its path — that one is the project it is about to create, not
// the one already open — and never for its quiet, which no interview wants.
func AnsweredForYou(command api.Command, field Field) bool {
	if !field.IsFlag {
		return false
	}
	if field.Id == quietFieldId && field.Type == typeBoolean {
		return true
	}
	return field.Id == pathFieldId && verbOf(command) != scaffoldVerb
}

// askField picks the question one field is asked as: yes or no for a boolean,
// a repeated question for an array, one question otherwise — each of them over
// a menu when the field names something that already exists.
func askField(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, field Field) ([]any, error) {
	if field.Type == typeBoolean {
		answer, err := sandbox.Deps.Interviewer.BoolQuestion(questionFor(sandbox, field))
		if err != nil {
			return nil, err
		}
		return []any{answer}, nil
	}

	offered := SuggestFor(sandbox, io, command, field.Id)

	if field.Array {
		return askArray(sandbox, field, offered)
	}

	return askScalar(sandbox, field, offered)
}

// askScalar asks for one value. A field with a candidate list is answered by
// choosing a row; an open list and an exhausted one both fall through to text.
func askScalar(sandbox *api.Sandbox, field Field, offered suggestion) ([]any, error) {
	if len(offered.Options) > 0 {
		chosen, err := sandbox.Deps.Interviewer.SingleAlternativeQuestion(
			questionFor(sandbox, field),
			menuRows(sandbox, offered, field),
		)
		if err != nil {
			return nil, err
		}

		switch chosen {
		case skipOptionId:
			return []any{}, nil
		case otherOptionId:
		default:
			value, ok := convert(sandbox, field, chosen)
			if ok {
				return []any{value}, nil
			}
		}
	}

	return askScalarText(sandbox, field)
}

// askScalarText reads one typed value, asking again until it converts and
// fits the bounds its declaration carries. An empty answer takes the default,
// unless the field is required, in which case there is nothing to fall back to.
func askScalarText(sandbox *api.Sandbox, field Field) ([]any, error) {
	for {
		answer, err := sandbox.Deps.Interviewer.StrQuestion(questionFor(sandbox, field))
		if err != nil {
			return nil, err
		}

		if sandbox.Deps.Stringsdeps.TrimSpace(answer) == "" {
			if !field.Required {
				return []any{}, nil
			}
			notice(sandbox, "%s has to be answered", label(field))
			continue
		}

		value, ok := convert(sandbox, field, answer)
		if !ok {
			notice(sandbox, "%q is not a %s", answer, typeWord(field.Type))
			continue
		}
		if !withinBounds(sandbox, field, value) {
			continue
		}

		return []any{value}, nil
	}
}

// askArray collects every occurrence of a repeatable field: by ticking rows
// when the field has a candidate list, and by asking for one value at a time
// until an empty answer when it does not.
func askArray(sandbox *api.Sandbox, field Field, offered suggestion) ([]any, error) {
	if len(offered.Options) > 0 && !offered.Open {
		chosen, err := sandbox.Deps.Interviewer.MultipleAlternativeQuestion(
			questionFor(sandbox, field),
			offered.Options,
		)
		if err != nil {
			return nil, err
		}

		values := []any{}
		for _, one := range chosen {
			if value, ok := convert(sandbox, field, one); ok {
				values = append(values, value)
			}
		}

		if len(values) == 0 && field.Required {
			notice(sandbox, "%s needs at least one answer", label(field))
			return askArray(sandbox, field, offered)
		}

		return values, nil
	}

	values := []any{}
	for {
		answer, err := sandbox.Deps.Interviewer.StrQuestion(arrayQuestionFor(sandbox, field, len(values)))
		if err != nil {
			return nil, err
		}

		if sandbox.Deps.Stringsdeps.TrimSpace(answer) == "" {
			if field.Required && len(values) == 0 {
				notice(sandbox, "%s needs at least one answer", label(field))
				continue
			}
			return values, nil
		}

		value, ok := convert(sandbox, field, answer)
		if !ok {
			notice(sandbox, "%q is not a %s", answer, typeWord(field.Type))
			continue
		}
		if !withinBounds(sandbox, field, value) {
			continue
		}

		values = append(values, value)
	}
}

// ─── Building the question ──────────────────────────────────────────────────

// menuRows is a candidate list with the rows the interview adds to it: one for
// leaving a field alone, and one for typing a value an open list does not
// hold. Both say what they do rather than what they are called — the menu is
// read by someone who has not seen the command before.
func menuRows(sandbox *api.Sandbox, offered suggestion, field Field) []interviewer.AlternativeOption {
	rows := []interviewer.AlternativeOption{}
	rows = append(rows, offered.Options...)

	if offered.Open {
		rows = append(rows, interviewer.AlternativeOption{Id: otherOptionId, Msg: "· none of these — let me type it"})
	}
	if !field.Required {
		rows = append(rows, interviewer.AlternativeOption{Id: skipOptionId, Msg: skipText(sandbox, field)})
	}

	return rows
}

// skipText is what the row leaving a field alone says will happen: the
// declared default is taken, or nothing is passed at all.
func skipText(sandbox *api.Sandbox, field Field) string {
	if field.HasDefault {
		return sandbox.Deps.Std.Sprintf("· skip it — %q is used", field.Default)
	}
	return "· skip it — leave it unset"
}

// questionFor words one field as a question: what it is called, what it is
// for, and what happens if it is left alone.
func questionFor(sandbox *api.Sandbox, field Field) string {
	question := label(field)
	if sandbox.Deps.Stringsdeps.TrimSpace(field.Description) != "" {
		question = sandbox.Deps.Std.Sprintf("%s — %s", question, field.Description)
	}
	return sandbox.Deps.Std.Sprintf("%s  (%s)", question, conditionOf(sandbox, field))
}

// arrayQuestionFor words the repeated question, counting what has been given
// so far so it is clear the answer is one of several.
func arrayQuestionFor(sandbox *api.Sandbox, field Field, given int) string {
	if given == 0 {
		return sandbox.Deps.Std.Sprintf("%s  (more than one allowed — answer nothing to stop)", questionFor(sandbox, field))
	}
	return sandbox.Deps.Std.Sprintf("%s — one more, or nothing to stop (%d so far)", label(field), given)
}

// conditionOf is the parenthesis after a question, in the words of someone who
// has not read a help screen: what kind of answer is expected, whether one is
// demanded, what stands in for it when it is skipped, and the bounds a number
// has to fall inside.
func conditionOf(sandbox *api.Sandbox, field Field) string {
	condition := typeWord(field.Type)

	if field.Required {
		condition += ", required"
	} else if field.HasDefault {
		condition = sandbox.Deps.Std.Sprintf("%s, %q if you skip it", condition, field.Default)
	} else {
		condition += ", optional"
	}

	if field.HasMin {
		condition = sandbox.Deps.Std.Sprintf("%s, %s or more", condition, numberLabel(sandbox, field.Type, field.Min))
	}
	if field.HasMax {
		condition = sandbox.Deps.Std.Sprintf("%s, %s or less", condition, numberLabel(sandbox, field.Type, field.Max))
	}

	return condition
}

// typeWord is a declared type said as the answer it asks for.
func typeWord(kind string) string {
	switch kind {
	case typeInt:
		return "whole number"
	case typeFloat:
		return "number"
	case typeBoolean:
		return "yes or no"
	}
	return "text"
}

// label is how a field is named back to the person: a flag by the spelling it
// answers to, an arg by the id it is declared under.
func label(field Field) string {
	if field.IsFlag && len(field.Identifiers) > 0 {
		return field.Identifiers[0]
	}
	return field.Id
}

// ─── Values ─────────────────────────────────────────────────────────────────

// convert reads one typed answer as the Go value its declaration names. The
// readers on api.Command type-assert, so a value bound under the wrong type
// reads back as the zero value and the handler silently does the wrong thing.
func convert(sandbox *api.Sandbox, field Field, raw string) (any, bool) {
	trimmed := sandbox.Deps.Stringsdeps.TrimSpace(raw)

	switch field.Type {
	case typeBoolean:
		return trimmed == "true", true
	case typeInt:
		value, err := sandbox.Deps.Stringsdeps.Atoi(trimmed)
		if err != nil {
			return nil, false
		}
		return value, true
	case typeFloat:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(trimmed, 64)
		if err != nil {
			return nil, false
		}
		return value, true
	}

	return raw, true
}

// DefaultValue is a field's declared default in the type the declaration
// names. entries.yaml spells every default as text, so the conversion the
// dispatch does when it binds one has to be done here too.
func DefaultValue(sandbox *api.Sandbox, field Field) any {
	value, ok := convert(sandbox, field, field.Default)
	if !ok {
		return field.Default
	}
	return value
}

// withinBounds enforces the min and max a numeric field declares, reporting
// the one it broke. A value of any other type carries no bounds and passes.
func withinBounds(sandbox *api.Sandbox, field Field, value any) bool {
	if !field.HasMin && !field.HasMax {
		return true
	}

	number, ok := numberOf(value)
	if !ok {
		return true
	}

	if field.HasMin && number < field.Min {
		notice(sandbox, "%s has to be %s or more", label(field), numberLabel(sandbox, field.Type, field.Min))
		return false
	}
	if field.HasMax && number > field.Max {
		notice(sandbox, "%s has to be %s or less", label(field), numberLabel(sandbox, field.Type, field.Max))
		return false
	}

	return true
}

// numberOf reads a bound value as the number a bound is compared against.
func numberOf(value any) (float64, bool) {
	if number, ok := value.(int); ok {
		return float64(number), true
	}
	if number, ok := value.(float64); ok {
		return number, true
	}
	return 0, false
}

// numberLabel spells a bound the way its declaration does.
func numberLabel(sandbox *api.Sandbox, kind string, value float64) string {
	if kind == typeInt {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}
