package triggerconf

// Trigger is the condition one declared slice or value has to meet for its
// unit — a route of route.yaml, a command of command.yaml — to run at all: the
// text compared, and how. Exists is false on an entry that declares none — a
// plain capture, or a value that is bound and never matched on.
type Trigger struct {
	Exists bool
	Type   string // "equal" | "prefix" | "text-prefix" | "suffix" | "regex" | "one-of"
	Value  string
	// Values are the texts a `one-of` trigger accepts; empty on every other
	// type.
	Values     []string
	Negate     bool
	IgnoreCase bool
}
