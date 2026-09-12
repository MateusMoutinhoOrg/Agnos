package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CommandIdentifier normalizes a user-typed command name into the CLI verb:
// lowercased, spaces and underscores turned into dashes
// ("My Feature" -> "my-feature").
func CommandIdentifier(sandbox *api.Sandbox, name string) string {
	out := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(name))
	out = sandbox.Deps.Stringsdeps.ReplaceAll(out, " ", "-")
	out = sandbox.Deps.Stringsdeps.ReplaceAll(out, "_", "-")
	return out
}

// ValidateCommandName reports whether a user-typed command name normalizes to
// a usable CLI verb. The identifier becomes a directory name and a Go package
// clause, so anything outside [a-z][a-z0-9-]* would be propagated straight
// into `package bad_name!` and break the whole project's build. The check runs
// before any file is written.
func ValidateCommandName(sandbox *api.Sandbox, name string) error {
	identifier := CommandIdentifier(sandbox, name)
	if identifier == "" {
		return sandbox.Deps.Std.Errorf("a command needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return sandbox.Deps.Std.Errorf("invalid command name %q: a command name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return sandbox.Deps.Std.Errorf(
				"invalid command name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory sandbox/internal/commands/%s and a Go package name)",
				name, CommandPackage(sandbox, name))
		}
	}
	return nil
}

// NoteNormalizedCommandName logs the rewriting a command name went through,
// so a name that is silently lowercased or re-punctuated ("MyCmd" -> "mycmd")
// is never a surprise later.
func NoteNormalizedCommandName(sandbox *api.Sandbox, name string) {
	identifier := CommandIdentifier(sandbox, name)
	if identifier != name {
		sandbox.Deps.Std.Log("note: command name %q normalized to %q \n", name, identifier)
	}
}

// CommandPackage is the Go package / directory name for a command: the
// identifier with dashes turned into underscores ("my-feature" -> "my_feature").
func CommandPackage(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.Stringsdeps.ReplaceAll(CommandIdentifier(sandbox, name), "-", "_")
}

// CommandDir is the project-relative directory holding a command package.
func CommandDir(sandbox *api.Sandbox, name string) string {
	return "sandbox/internal/commands/" + CommandPackage(sandbox, name)
}

// CommandEntriesPath is the project-relative path of a command's entries.yaml.
func CommandEntriesPath(sandbox *api.Sandbox, name string) string {
	return CommandDir(sandbox, name) + "/entries.yaml"
}

// LoadCommandConf reads and parses sandbox/internal/commands/<name>/entries.yaml.
func LoadCommandConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string) (*commandconf.CommandConf, error) {
	if err := ValidateCommandName(sandbox, name); err != nil {
		return nil, err
	}
	content, err := io.ReadFile(CommandEntriesPath(sandbox, name))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("command %q not found in %s", CommandIdentifier(sandbox, name), CommandDir(sandbox, name))
	}
	conf, err := commandconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("commands/%s/entries.yaml: %w", CommandPackage(sandbox, name), err)
	}
	return conf, nil
}

// SaveCommandConf renders conf back over sandbox/internal/commands/<name>/entries.yaml.
func SaveCommandConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string, conf *commandconf.CommandConf) error {
	return io.WriteFileOverwrite(CommandEntriesPath(sandbox, name), []byte(conf.Render()))
}

// FieldName normalizes a flag/arg name the same way command names are
// ("Out File" -> "out-file"); the generated Go field is derived from it.
func FieldName(sandbox *api.Sandbox, name string) string {
	return CommandIdentifier(sandbox, name)
}

// NewField builds a commandconf.Field from the raw values typed on the
// command line, validating the type and parsing the default/min/max
// literals. Identifiers are left exactly as given (the caller decides
// whether the field is a flag or a positional arg).
func NewField(sandbox *api.Sandbox, props api.FieldProps) (commandconf.Field, error) {
	field := commandconf.Field{
		Key:         FieldName(sandbox, props.Name),
		Identifiers: props.Identifiers,
		Description: sandbox.Deps.Stringsdeps.TrimSpace(props.Description),
		Examples:    props.Examples,
		Required:    props.Required,
		Array:       props.Array,
	}
	if field.Key == "" {
		return field, sandbox.Deps.Std.Errorf("a flag/arg needs a name")
	}

	kind, ok := FieldType(sandbox, props.Type)
	if !ok {
		return field, sandbox.Deps.Std.Errorf("unknown type %q (use string, boolean, int or float)", props.Type)
	}
	field.Type = kind

	if field.Required && kind == "boolean" {
		return field, sandbox.Deps.Std.Errorf("a boolean field cannot be required (its absence already means false)")
	}

	if props.Default != "" {
		if field.Required {
			return field, sandbox.Deps.Std.Errorf("a field cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := checkLiteral(sandbox, kind, "default", props.Default); err != nil {
			return field, err
		}
		field.HasDefault = true
		field.Default = props.Default
	}

	if props.Min != "" {
		value, err := parseBound(sandbox, kind, "min", props.Min)
		if err != nil {
			return field, err
		}
		field.Min, field.HasMin = value, true
	}
	if props.Max != "" {
		value, err := parseBound(sandbox, kind, "max", props.Max)
		if err != nil {
			return field, err
		}
		field.Max, field.HasMax = value, true
	}
	if field.HasMin && field.HasMax && field.Min > field.Max {
		return field, sandbox.Deps.Std.Errorf("min (%s) is greater than max (%s)", props.Min, props.Max)
	}

	return field, nil
}

// FieldType maps the type spellings accepted on the command line onto the
// canonical entries.yaml set; "" defaults to string.
func FieldType(sandbox *api.Sandbox, raw string) (string, bool) {
	switch sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw)) {
	case "", "string", "str":
		return "string", true
	case "bool", "boolean":
		return "boolean", true
	case "int", "integer":
		return "int", true
	case "float", "double", "number":
		return "float", true
	default:
		return "", false
	}
}

// FindField returns the index of the field named name in fields, or -1.
func FindField(sandbox *api.Sandbox, fields []commandconf.Field, name string) int {
	key := FieldName(sandbox, name)
	for i, field := range fields {
		if field.Key == key {
			return i
		}
	}
	return -1
}

// AppendPosition is the --position value meaning "put it last". It is the
// flag's default, so "not given" and "at the end" are the same request.
const AppendPosition = -1

// CheckPosition validates a --position against the list it will be inserted
// into: AppendPosition means the end, and anything else must land inside
// 0..len(fields). An out-of-range index is reported as such instead of
// silently becoming an append (which then trips a later, unrelated rule).
func CheckPosition(sandbox *api.Sandbox, kind string, position int, fields []commandconf.Field) (int, error) {
	if position == AppendPosition {
		return len(fields), nil
	}
	if position < 0 {
		return 0, sandbox.Deps.Std.Errorf("--position %d is negative: use an index from 0 to %d, or leave it out to append", position, len(fields))
	}
	if position > len(fields) {
		return 0, sandbox.Deps.Std.Errorf("--position %d is out of range: this command has %d %s(s), so the accepted range is 0 to %d", position, len(fields), kind, len(fields))
	}
	return position, nil
}

// InsertField places field at position inside fields (appending when
// position is negative or past the end).
func InsertField(fields []commandconf.Field, field commandconf.Field, position int) []commandconf.Field {
	if position < 0 || position >= len(fields) {
		return append(fields, field)
	}
	out := make([]commandconf.Field, 0, len(fields)+1)
	out = append(out, fields[:position]...)
	out = append(out, field)
	out = append(out, fields[position:]...)
	return out
}

// RemoveField drops the field at index from fields.
func RemoveField(fields []commandconf.Field, index int) []commandconf.Field {
	out := make([]commandconf.Field, 0, len(fields)-1)
	out = append(out, fields[:index]...)
	out = append(out, fields[index+1:]...)
	return out
}

// AppendUnique appends each value of extra to values, skipping duplicates.
func AppendUnique(values []string, extra []string) []string {
	for _, candidate := range extra {
		found := false
		for _, existing := range values {
			if existing == candidate {
				found = true
				break
			}
		}
		if !found {
			values = append(values, candidate)
		}
	}
	return values
}

func checkLiteral(sandbox *api.Sandbox, kind string, label string, raw string) error {
	switch kind {
	case "boolean":
		if raw != "true" && raw != "false" {
			return sandbox.Deps.Std.Errorf("%s for a boolean must be true or false, got %q", label, raw)
		}
	case "int":
		if _, err := sandbox.Deps.Stringsdeps.ParseInt(raw, 10, 64); err != nil {
			return sandbox.Deps.Std.Errorf("%s must be an int, got %q", label, raw)
		}
	case "float":
		if _, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64); err != nil {
			return sandbox.Deps.Std.Errorf("%s must be a float, got %q", label, raw)
		}
	}
	return nil
}

func parseBound(sandbox *api.Sandbox, kind string, label string, raw string) (float64, error) {
	if kind != "int" && kind != "float" {
		return 0, sandbox.Deps.Std.Errorf("%s only applies to int/float fields", label)
	}
	if err := checkLiteral(sandbox, kind, label, raw); err != nil {
		return 0, err
	}
	value, _ := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
	return value, nil
}
