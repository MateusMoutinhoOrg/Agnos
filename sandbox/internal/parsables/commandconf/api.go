package commandconf

// Field is one flag or one positional argument declared in a command's
// entries.yaml. Flags carry Identifiers ("--path", "-p"); positional args
// leave Identifiers empty and are matched by order, Key naming the generated
// struct field.
type Field struct {
	Key         string
	Identifiers []string
	Description string
	Examples    []string
	Type        string // "string" | "boolean" | "int" | "float"
	Default     string
	HasDefault  bool
	Required    bool
	Array       bool
	Min         float64
	HasMin      bool
	Max         float64
	HasMax      bool
}

// CommandConf is the parsed form of sandbox/internal/commands/<name>/entries.yaml
// — the declarative description of one command the user writes by hand and
// `agnos build` turns into an entries.go struct plus a dispatch arm in
// sandbox/internal/cli/climain.go.
type CommandConf struct {
	Identifiers     []string
	Category        string
	Help            string
	LongDescription string
	Examples        []string
	Hidden          bool
	Flags           []Field
	Args            []Field

	Render func() string
}
