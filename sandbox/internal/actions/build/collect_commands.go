package build

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CollectCommands reads every sandbox/internal/commands/<name>/entries.yaml and
// returns one rich data map per command, in listing order, for the
// {{range .Commands}} loops in the generated sandbox/internal/cli/climain.go and
// help package. `help` is collected like every other command: its entries.yaml
// is written by GenerateHelpEntriesYaml just before this runs.
func CollectCommands(deps *deps.Deps, io *smartio.SmartIO) ([]map[string]any, error) {
	var commands []map[string]any

	for _, dir := range io.ListDirs("sandbox/internal/commands") {
		name := lastSegmentOf(dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile("sandbox/internal/commands/" + name + "/entries.yaml")
		if err != nil {
			continue
		}

		conf, err := commandconf.New(deps, string(content))
		if err != nil {
			return nil, deps.Std.Errorf("commands/%s/entries.yaml: %w", name, err)
		}

		commands = append(commands, commandData(name, conf))
	}

	return commands, nil
}

func commandData(name string, conf *commandconf.CommandConf) map[string]any {
	flags := make([]map[string]any, 0, len(conf.Flags))
	for _, flag := range conf.Flags {
		flags = append(flags, fieldData(flag))
	}
	args := make([]map[string]any, 0, len(conf.Args))
	for _, arg := range conf.Args {
		args = append(args, fieldData(arg))
	}

	return map[string]any{
		"Name":            name,
		"QuietField":      quietField(conf),
		"GoName":          exportedName(name),
		"Identifiers":     conf.Identifiers,
		"IdentifiersGo":   goStringList(conf.Identifiers),
		"MatchExpr":       matchExpr(conf.Identifiers),
		"Category":        conf.Category,
		"Help":            conf.Help,
		"LongDescription": conf.LongDescription,
		"Examples":        conf.Examples,
		"Hidden":          conf.Hidden,
		"Flags":           flags,
		"Args":            args,
	}
}

func fieldData(field commandconf.Field) map[string]any {
	return map[string]any{
		"Key":            field.Key,
		"GoField":        exportedName(field.Key),
		"GoType":         goType(field.Type, field.Array),
		"IsBool":         field.Type == "boolean",
		"IsArray":        field.Array,
		"Identifiers":    field.Identifiers,
		"IdentifiersGo":  goStringList(field.Identifiers),
		"OptionGetter":   optionGetter(field.Type),
		"ArgGetter":      argGetter(field.Type),
		"ParseFunc":      parseFunc(field.Type),
		"Required":       field.Required,
		"HasDefault":     field.HasDefault,
		"DefaultLiteral": defaultLiteral(field.Type, field.Default),
		"ElemLiteral":    elemLiteral(field.Type),
		"Description":    field.Description,
		"Examples":       field.Examples,
		"Type":           field.Type,
		"Default":        field.Default,
		"MinLabel":       numberLabel(field, field.Min, field.HasMin),
		"MaxLabel":       numberLabel(field, field.Max, field.HasMax),
		"RangeCheck":     rangeCheck(field),
	}
}

// numberLabel renders a min/max bound as the literal it has in entries.yaml
// ("" when the bound is unset), for help display.
func numberLabel(field commandconf.Field, value float64, has bool) string {
	if !has {
		return ""
	}
	if field.Type == "int" {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'g', -1, 64)
}

// rangeCheck emits the Go statements the generated dispatch runs, after a
// numeric flag/arg has been bound, to enforce its min/max bounds. It returns
// "" for fields that carry no bound (or are not int/float scalars). The body
// is indented two tabs — the depth of the block it is spliced into.
func rangeCheck(field commandconf.Field) string {
	if field.Array || (field.Type != "int" && field.Type != "float") {
		return ""
	}
	if !field.HasMin && !field.HasMax {
		return ""
	}

	subject := "flag"
	if len(field.Identifiers) == 0 {
		subject = "arg"
	}
	goField := exportedName(field.Key)

	var b strings.Builder
	// failOp is the comparison that means "out of range"; wantOp is what the
	// message tells the user to satisfy.
	guard := func(failOp, wantOp, bound string) {
		fmt.Fprintf(&b,
			"\t\tif entries.%s %s %s {\n"+
				"\t\t\tdeps.Std.Error(\"%s '%s' must be %s %s\\n\")\n"+
				"\t\t\treturn ExitUsage\n"+
				"\t\t}\n",
			goField, failOp, bound, subject, field.Key, wantOp, bound)
	}
	if field.HasMin {
		guard("<", ">=", numberLabel(field, field.Min, true))
	}
	if field.HasMax {
		guard(">", "<=", numberLabel(field, field.Max, true))
	}
	return strings.TrimRight(b.String(), "\n")
}

// ─── helpers ────────────────────────────────────────────────────────────────

func lastSegmentOf(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// exportedName turns a kebab/snake identifier into an exported Go name:
// "project-name" -> "ProjectName", "unsafe" -> "Unsafe".
func exportedName(raw string) string {
	parts := strings.FieldsFunc(raw, func(r rune) bool { return r == '-' || r == '_' })
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	if b.Len() == 0 {
		return "Field"
	}
	return b.String()
}

func goType(kind string, array bool) string {
	base := "string"
	switch kind {
	case "boolean":
		base = "bool"
	case "int":
		base = "int"
	case "float":
		base = "float64"
	}
	if array {
		return "[]" + base
	}
	return base
}

// elemLiteral is the zero literal for one element of a value flag/arg, used
// when appending to an array field.
func elemLiteral(kind string) string {
	switch kind {
	case "int":
		return "0"
	case "float":
		return "0"
	default:
		return `""`
	}
}

// quietField is the generated Entries field the dispatch checks to turn the
// progress channel off, "" when the command declares no boolean quiet flag.
// Every command that wants --quiet to work declares the flag; the dispatch
// does the silencing once, so no handler has to.
func quietField(conf *commandconf.CommandConf) string {
	for _, flag := range conf.Flags {
		if flag.Key == "quiet" && flag.Type == "boolean" && !flag.Array {
			return exportedName(flag.Key)
		}
	}
	return ""
}

// parseFunc is the generated helper that converts one raw command-line value
// into the field's Go type, reporting a clean usage error of its own.
func parseFunc(kind string) string {
	switch kind {
	case "int":
		return "parseIntValue"
	case "float":
		return "parseFloatValue"
	default:
		return "parseStringValue"
	}
}

func optionGetter(kind string) string {
	switch kind {
	case "int":
		return "IntOption"
	case "float":
		return "DoubleOption"
	default:
		return "StringOption"
	}
}

func argGetter(kind string) string {
	switch kind {
	case "int":
		return "Int"
	case "float":
		return "Double"
	default:
		return "String"
	}
}

func defaultLiteral(kind, value string) string {
	switch kind {
	case "boolean":
		if value == "true" {
			return "true"
		}
		return "false"
	case "int":
		if value == "" {
			return "0"
		}
		return value
	case "float":
		if value == "" {
			return "0"
		}
		return value
	default:
		return strconv.Quote(value)
	}
}

func goStringList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, strconv.Quote(value))
	}
	return strings.Join(quoted, ", ")
}

// matchExpr builds the boolean switch guard that matches the command verb,
// e.g. `action == "version" || action == "--version"`.
func matchExpr(identifiers []string) string {
	if len(identifiers) == 0 {
		return "false"
	}
	parts := make([]string, 0, len(identifiers))
	for _, id := range identifiers {
		parts = append(parts, "action == "+strconv.Quote(id))
	}
	return strings.Join(parts, " || ")
}
