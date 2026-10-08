package commandconf

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"

// Arg is one entry of a command's `args`: the segments of the command line from
// Start to End, both inclusive, End -1 standing for the last one. One segment
// is bound to Input.<Id> in its Type; a range as a []string. A segment is a
// leading word of the command line — every token before the first one starting
// with "-" — or any token after a bare "--".
type Arg struct {
	Id          string
	Start       int
	End         int
	Type        string // "string" | "integer" | "number" | "uuid"
	Required    bool
	Default     string
	HasDefault  bool
	Description string
	Trigger     triggerconf.Trigger
}

// Flag is one entry of a command's `flags`: the value following one of Keys on
// the command line, converted to Type and bound to Input.<Id>.
type Flag struct {
	Id string
	// Keys are the spellings a user types ("--command", "-c"); HasKeys is
	// false on a flag declaring none, which answers to --<id in kebab-case>.
	Keys        []string
	HasKeys     bool
	Type        string // "string" | "integer" | "number" | "boolean" | "string-array" | "integer-array"
	Required    bool
	Default     string
	HasDefault  bool
	Min         float64
	HasMin      bool
	Max         float64
	HasMax      bool
	Enum        []string
	Pattern     string
	Description string
	Trigger     triggerconf.Trigger
}

// CommandConf is the parsed form of sandbox/internal/commands/<name>/command.yaml
// — the declaration of one command, which `agnos build` turns into the
// api.Command of that command's generated new.go and the Input of its
// generated input.go.
type CommandConf struct {
	// Priority is the rung the command runs on; HasPriority is false on a
	// declaration missing it, which verify reports.
	Priority    int
	HasPriority bool
	// Segments is how many segments the command line has to have; HasSegments
	// is false on a command that takes any count.
	Segments    int
	HasSegments bool
	// Strict reports that every token of the command line has to be read by
	// the command or one run before it; a middleware declares false.
	Strict      bool
	Args        []Arg
	Flags       []Flag
	Category    string
	Summary     string
	Description string
	Examples    []string
	Hidden      bool
	// Legacy holds the keys of the declaration entries.yaml was, which verify
	// names as an old declaration.
	Legacy []string

	Render      func() string
	Pattern     func() string
	Identifiers func() []string
}
