package utils

import (
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CommandIdentifier normalizes a user-typed command name into the CLI verb:
// lowercased, spaces and underscores turned into dashes
// ("My Feature" -> "my-feature").
func CommandIdentifier(name string) string {
	out := strings.ToLower(strings.TrimSpace(name))
	out = strings.ReplaceAll(out, " ", "-")
	out = strings.ReplaceAll(out, "_", "-")
	return out
}

// ValidateCommandName reports whether a user-typed command name normalizes to
// a usable CLI verb. The identifier becomes a directory name and a Go package
// clause, so anything outside [a-z][a-z0-9-]* would be propagated straight
// into `package bad_name!` and break the whole project's build. The check runs
// before any file is written.
func ValidateCommandName(deps *deps.Deps, name string) error {
	identifier := CommandIdentifier(name)
	if identifier == "" {
		return deps.Std.Errorf("a command needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return deps.Std.Errorf("invalid command name %q: a command name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return deps.Std.Errorf(
				"invalid command name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory sandbox/internal/commands/%s and a Go package name)",
				name, CommandPackage(name))
		}
	}
	return nil
}

// NoteNormalizedCommandName logs the rewriting a command name went through,
// so a name that is silently lowercased or re-punctuated ("MyCmd" -> "mycmd")
// is never a surprise later.
func NoteNormalizedCommandName(deps *deps.Deps, name string) {
	identifier := CommandIdentifier(name)
	if identifier != name {
		deps.Std.Log("note: command name %q normalized to %q \n", name, identifier)
	}
}

// CommandPackage is the Go package / directory name for a command: the
// identifier with dashes turned into underscores ("my-feature" -> "my_feature").
func CommandPackage(name string) string {
	return strings.ReplaceAll(CommandIdentifier(name), "-", "_")
}

// CommandDir is the project-relative directory holding a command package.
func CommandDir(name string) string {
	return "sandbox/internal/commands/" + CommandPackage(name)
}

// CommandEntriesPath is the project-relative path of a command's entries.yaml.
func CommandEntriesPath(name string) string {
	return CommandDir(name) + "/entries.yaml"
}

// LoadCommandConf reads and parses sandbox/internal/commands/<name>/entries.yaml.
func LoadCommandConf(deps *deps.Deps, io *smartio.SmartIO, name string) (*commandconf.CommandConf, error) {
	if err := ValidateCommandName(deps, name); err != nil {
		return nil, err
	}
	content, err := io.ReadFile(CommandEntriesPath(name))
	if err != nil {
		return nil, deps.Std.Errorf("command %q not found in %s", CommandIdentifier(name), CommandDir(name))
	}
	conf, err := commandconf.New(deps, string(content))
	if err != nil {
		return nil, deps.Std.Errorf("commands/%s/entries.yaml: %w", CommandPackage(name), err)
	}
	return conf, nil
}

// SaveCommandConf renders conf back over sandbox/internal/commands/<name>/entries.yaml.
func SaveCommandConf(deps *deps.Deps, io *smartio.SmartIO, name string, conf *commandconf.CommandConf) error {
	return io.WriteFileOverwrite(CommandEntriesPath(name), []byte(conf.Render()))
}

// FieldName normalizes a flag/arg name the same way command names are
// ("Out File" -> "out-file"); the generated Go field is derived from it.
func FieldName(name string) string {
	return CommandIdentifier(name)
}

// NewField builds a commandconf.Field from the raw values typed on the
// command line, validating the type and parsing the default/min/max
// literals. Identifiers are left exactly as given (the caller decides
// whether the field is a flag or a positional arg).
func NewField(deps *deps.Deps, props api.FieldProps) (commandconf.Field, error) {
	field := commandconf.Field{
		Key:         FieldName(props.Name),
		Identifiers: props.Identifiers,
		Description: strings.TrimSpace(props.Description),
		Examples:    props.Examples,
		Required:    props.Required,
		Array:       props.Array,
	}
	if field.Key == "" {
		return field, deps.Std.Errorf("a flag/arg needs a name")
	}

	kind, ok := FieldType(props.Type)
	if !ok {
		return field, deps.Std.Errorf("unknown type %q (use string, boolean, int or float)", props.Type)
	}
	field.Type = kind

	if field.Required && kind == "boolean" {
		return field, deps.Std.Errorf("a boolean field cannot be required (its absence already means false)")
	}

	if props.Default != "" {
		if field.Required {
			return field, deps.Std.Errorf("a field cannot be both required and carry a default (the default already covers its absence)")
		}
		if err := checkLiteral(deps, kind, "default", props.Default); err != nil {
			return field, err
		}
		field.HasDefault = true
		field.Default = props.Default
	}

	if props.Min != "" {
		value, err := parseBound(deps, kind, "min", props.Min)
		if err != nil {
			return field, err
		}
		field.Min, field.HasMin = value, true
	}
	if props.Max != "" {
		value, err := parseBound(deps, kind, "max", props.Max)
		if err != nil {
			return field, err
		}
		field.Max, field.HasMax = value, true
	}
	if field.HasMin && field.HasMax && field.Min > field.Max {
		return field, deps.Std.Errorf("min (%s) is greater than max (%s)", props.Min, props.Max)
	}

	return field, nil
}

// FieldType maps the type spellings accepted on the command line onto the
// canonical entries.yaml set; "" defaults to string.
func FieldType(raw string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
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
func FindField(fields []commandconf.Field, name string) int {
	key := FieldName(name)
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
func CheckPosition(deps *deps.Deps, kind string, position int, fields []commandconf.Field) (int, error) {
	if position == AppendPosition {
		return len(fields), nil
	}
	if position < 0 {
		return 0, deps.Std.Errorf("--position %d is negative: use an index from 0 to %d, or leave it out to append", position, len(fields))
	}
	if position > len(fields) {
		return 0, deps.Std.Errorf("--position %d is out of range: this command has %d %s(s), so the accepted range is 0 to %d", position, len(fields), kind, len(fields))
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

func checkLiteral(deps *deps.Deps, kind string, label string, raw string) error {
	switch kind {
	case "boolean":
		if raw != "true" && raw != "false" {
			return deps.Std.Errorf("%s for a boolean must be true or false, got %q", label, raw)
		}
	case "int":
		if _, err := strconv.ParseInt(raw, 10, 64); err != nil {
			return deps.Std.Errorf("%s must be an int, got %q", label, raw)
		}
	case "float":
		if _, err := strconv.ParseFloat(raw, 64); err != nil {
			return deps.Std.Errorf("%s must be a float, got %q", label, raw)
		}
	}
	return nil
}

func parseBound(deps *deps.Deps, kind string, label string, raw string) (float64, error) {
	if kind != "int" && kind != "float" {
		return 0, deps.Std.Errorf("%s only applies to int/float fields", label)
	}
	if err := checkLiteral(deps, kind, label, raw); err != nil {
		return 0, err
	}
	value, _ := strconv.ParseFloat(raw, 64)
	return value, nil
}
