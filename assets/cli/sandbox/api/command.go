package api

// CommandArg is one positional argument a command declares — the parsed form
// of one entry under `args:` in that command's entries.yaml. It is matched by
// position, so Id names it for GetItem and for the help screen, never on the
// command line.
type CommandArg struct {
	// Type is the declared type: "string", "boolean", "int" or "float".
	Type string
	// Id is the name the argument is declared and read back under.
	Id string
	// Required reports that the command line is rejected without it.
	Required bool
	// Array reports that it takes every argument left on the line.
	Array bool
	// Description is the one-line help text.
	Description string
	// Examples are whole command lines the help screen prints.
	Examples []string
	// Default is the value bound when the argument is absent, spelled as
	// it is written in entries.yaml.
	Default string
	// HasDefault tells a declared empty default from no default at all.
	HasDefault bool
	// Min and Max bound a numeric value; HasMin and HasMax tell a bound of
	// zero from no bound.
	Min    float64
	HasMin bool
	Max    float64
	HasMax bool
}

// CommandFlag is one flag a command declares — the parsed form of one entry
// under `flags:` in that command's entries.yaml. It is matched by Identifiers
// and read back under Id.
type CommandFlag struct {
	// Type is the declared type: "string", "boolean", "int" or "float".
	Type string
	// Id is the name the flag is declared and read back under.
	Id string
	// Required reports that the command line is rejected without it.
	Required bool
	// Array reports that every occurrence is kept, not just the first.
	Array bool
	// Description is the one-line help text.
	Description string
	// Examples are whole command lines the help screen prints.
	Examples []string
	// Default is the value bound when the flag is absent, spelled as it is
	// written in entries.yaml.
	Default string
	// HasDefault tells a declared empty default from no default at all.
	HasDefault bool
	// Identifiers are the spellings the flag answers to ("--path", "-p").
	Identifiers []string
	// Min and Max bound a numeric value; HasMin and HasMax tell a bound of
	// zero from no bound.
	Min    float64
	HasMin bool
	Max    float64
	HasMax bool
}

// Command is one command of the project, as the sandbox offers it: the whole
// of what its entries.yaml declares, plus the handler behind it. Sandbox.Commands
// holds one per sandbox/internal/commands/<name>/, each built by that package's
// generated NewCommand. The cli dispatch fills Items from the command line and
// then calls Handler; a caller holding the sandbox can do the same.
type Command struct {
	// Name is the package directory of the command, snake_case.
	Name string
	// Identifiers are the verbs it answers to ("add-flag"), the first of
	// which is the one the help screen prints.
	Identifiers []string
	// Category groups it on the general help screen.
	Category string
	// Help is the one-line description.
	Help string
	// LongDescription is the paragraph the per-command help screen prints.
	LongDescription string
	// Examples are whole command lines the help screen prints.
	Examples []string
	// Hidden keeps it off the general help screen without disabling it.
	Hidden bool
	// Args are the positional arguments, in declaration order.
	Args []CommandArg
	// Flags are the flags, in declaration order.
	Flags []CommandFlag

	// Items holds the values bound to the declaration: one entry per flag
	// or arg Id, in declaration order for an array, one element long for a
	// scalar, absent when nothing was bound. The Get* fields read it.
	Items map[string][]any

	// GetItem returns every value bound under one Id, nil when none was.
	GetItem func(id string) []any
	// GetString returns the first string bound under one Id, "" when none
	// was.
	GetString func(id string) string
	// GetBool returns the first boolean bound under one Id, false when none
	// was.
	GetBool func(id string) bool
	// GetInt returns the first int bound under one Id, 0 when none was.
	GetInt func(id string) int
	// GetFloat returns the first float bound under one Id, 0 when none was.
	GetFloat func(id string) float64
	// GetStrings returns every string bound under one Id, in order.
	GetStrings func(id string) []string
	// GetInts returns every int bound under one Id, in order.
	GetInts func(id string) []int
	// GetFloats returns every float bound under one Id, in order.
	GetFloats func(id string) []float64

	// Handler runs the command against the values in Items. It is the
	// command package's own CommandHandler, closed over the sandbox.
	Handler func() int
}

// NewCommand returns an empty Command with Items open and every Get* reader
// bound to it. A generated NewCommand fills the declaration and the handler on
// top of what this returns, so every command reads its values the same way.
func NewCommand() *Command {
	command := &Command{
		Identifiers: []string{},
		Examples:    []string{},
		Args:        []CommandArg{},
		Flags:       []CommandFlag{},
		Items:       map[string][]any{},
	}

	command.GetItem = func(id string) []any {
		return command.Items[id]
	}
	command.GetString = func(id string) string {
		value, _ := first(command, id).(string)
		return value
	}
	command.GetBool = func(id string) bool {
		value, _ := first(command, id).(bool)
		return value
	}
	command.GetInt = func(id string) int {
		value, _ := first(command, id).(int)
		return value
	}
	command.GetFloat = func(id string) float64 {
		value, _ := first(command, id).(float64)
		return value
	}
	command.GetStrings = func(id string) []string {
		values := []string{}
		for _, item := range command.Items[id] {
			if value, ok := item.(string); ok {
				values = append(values, value)
			}
		}
		return values
	}
	command.GetInts = func(id string) []int {
		values := []int{}
		for _, item := range command.Items[id] {
			if value, ok := item.(int); ok {
				values = append(values, value)
			}
		}
		return values
	}
	command.GetFloats = func(id string) []float64 {
		values := []float64{}
		for _, item := range command.Items[id] {
			if value, ok := item.(float64); ok {
				values = append(values, value)
			}
		}
		return values
	}

	return command
}

// first is the one value of an Id a scalar reader wants, nil when nothing was
// bound under it.
func first(command *Command, id string) any {
	items := command.Items[id]
	if len(items) == 0 {
		return nil
	}
	return items[0]
}
