package api

// TriggerType is how a Trigger compares the text it is handed.
type TriggerType int

const (
	// EqualTrigger matches a text that is exactly the trigger's Value.
	EqualTrigger TriggerType = iota
	// PrefixTrigger matches a text that is the Value or continues it with a
	// new segment: "/admin" matches "/admin" and "/admin/users", never
	// "/administrator". A Value of "/" matches every path.
	PrefixTrigger
	// TextPrefixTrigger matches a text that begins with the Value, whatever
	// follows it: "/admin" matches "/administrator" too.
	TextPrefixTrigger
	// SuffixTrigger matches a text that ends with the Value.
	SuffixTrigger
	// RegexTrigger matches a text the Value, a regular expression, matches.
	RegexTrigger
	// OneOfTrigger matches a text that is exactly one of the trigger's
	// Values — how a command answers to more than one name.
	OneOfTrigger
)

// Trigger is the condition one declared slice or value has to meet for its
// unit to run at all — the parsed form of one `trigger:` of a route.yaml or a
// command.yaml. Failing it is a non-match, never a usage error or a 400: the
// input is for some other unit.
type Trigger struct {
	// Exist tells a declared trigger from none at all; an entry with none
	// matches whatever the request brought.
	Exist bool
	// Type is how Value is compared.
	Type TriggerType
	// Value is what the text is compared against.
	Value string
	// Values are the texts a OneOfTrigger accepts; nil on every other type.
	Values []string
	// Negate inverts the comparison: the trigger holds when the text does
	// not match.
	Negate bool
	// IgnoreCase compares without regard to case.
	IgnoreCase bool
}
