package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectCommands reads every sandbox/internal/commands/<name>/entries.yaml and
// returns one rich data map per command, in listing order, for the
// {{range .Commands}} loops in the generated sandbox/internal/cli/climain.go and
// help package. `help` is collected like every other command: its entries.yaml
// is written by GenerateHelpEntriesYaml just before this runs.
func CollectCommands(sandbox *api.Sandbox, io *smartio.SmartIO) ([]map[string]any, error) {
	var commands []map[string]any

	for _, dir := range io.ListDirs("sandbox/internal/commands") {
		name := lastSegmentOf(sandbox, dir)
		if name == "" {
			continue
		}

		content, err := io.ReadFile("sandbox/internal/commands/" + name + "/entries.yaml")
		if err != nil {
			continue
		}

		conf, err := commandconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("commands/%s/entries.yaml: %w", name, err)
		}

		commands = append(commands, commandData(sandbox, name, conf))
	}

	return commands, nil
}

func commandData(sandbox *api.Sandbox, name string, conf *commandconf.CommandConf) map[string]any {
	flags := make([]map[string]any, 0, len(conf.Flags))
	for _, flag := range conf.Flags {
		flags = append(flags, fieldData(sandbox, flag))
	}
	args := make([]map[string]any, 0, len(conf.Args))
	for _, arg := range conf.Args {
		args = append(args, fieldData(sandbox, arg))
	}

	return map[string]any{
		"Name":            name,
		"QuietField":      quietField(sandbox, conf),
		"GoName":          exportedName(sandbox, name),
		"Identifiers":     conf.Identifiers,
		"IdentifiersGo":   goStringList(sandbox, conf.Identifiers),
		"MatchExpr":       matchExpr(sandbox, conf.Identifiers),
		"Category":        conf.Category,
		"Help":            conf.Help,
		"LongDescription": conf.LongDescription,
		"Examples":        conf.Examples,
		"Hidden":          conf.Hidden,
		"Flags":           flags,
		"Args":            args,
	}
}

func fieldData(sandbox *api.Sandbox, field commandconf.Field) map[string]any {
	return map[string]any{
		"Key":            field.Key,
		"GoField":        exportedName(sandbox, field.Key),
		"GoType":         goType(field.Type, field.Array),
		"IsBool":         field.Type == "boolean",
		"IsArray":        field.Array,
		"Identifiers":    field.Identifiers,
		"IdentifiersGo":  goStringList(sandbox, field.Identifiers),
		"OptionGetter":   optionGetter(field.Type),
		"ArgGetter":      argGetter(field.Type),
		"ParseFunc":      parseFunc(field.Type),
		"Required":       field.Required,
		"HasDefault":     field.HasDefault,
		"DefaultLiteral": defaultLiteral(sandbox, field.Type, field.Default),
		"ElemLiteral":    elemLiteral(field.Type),
		"Description":    field.Description,
		"Examples":       field.Examples,
		"Type":           field.Type,
		"Default":        field.Default,
		"MinLabel":       numberLabel(sandbox, field, field.Min, field.HasMin),
		"MaxLabel":       numberLabel(sandbox, field, field.Max, field.HasMax),
		"RangeCheck":     rangeCheck(sandbox, field),
	}
}

// numberLabel renders a min/max bound as the literal it has in entries.yaml
// ("" when the bound is unset), for help display.
func numberLabel(sandbox *api.Sandbox, field commandconf.Field, value float64, has bool) string {
	if !has {
		return ""
	}
	if field.Type == "int" {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}

// rangeCheck emits the Go statements the generated dispatch runs, after a
// numeric flag/arg has been bound, to enforce its min/max bounds. It returns
// "" for fields that carry no bound (or are not int/float scalars). The body
// is indented two tabs — the depth of the block it is spliced into.
func rangeCheck(sandbox *api.Sandbox, field commandconf.Field) string {
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
	goField := exportedName(sandbox, field.Key)

	b := ""
	// failOp is the comparison that means "out of range"; wantOp is what the
	// message tells the user to satisfy.
	guard := func(failOp, wantOp, bound string) {
		b += sandbox.Deps.Std.Sprintf(
			"\t\tif entries.%s %s %s {\n"+
				"\t\t\tsandbox.Deps.Std.Error(\"%s '%s' must be %s %s\\n\")\n"+
				"\t\t\treturn ExitUsage\n"+
				"\t\t}\n",
			goField, failOp, bound, subject, field.Key, wantOp, bound)
	}
	if field.HasMin {
		guard("<", ">=", numberLabel(sandbox, field, field.Min, true))
	}
	if field.HasMax {
		guard(">", "<=", numberLabel(sandbox, field, field.Max, true))
	}
	return sandbox.Deps.Stringsdeps.TrimRight(b, "\n")
}

// ─── helpers ────────────────────────────────────────────────────────────────

func lastSegmentOf(sandbox *api.Sandbox, path string) string {
	parts := sandbox.Deps.Stringsdeps.Split(path, "/")
	return parts[len(parts)-1]
}

// exportedName turns a kebab/snake identifier into an exported Go name:
// "project-name" -> "ProjectName", "unsafe" -> "Unsafe".
func exportedName(sandbox *api.Sandbox, raw string) string {
	parts := sandbox.Deps.Stringsdeps.FieldsFunc(raw, func(r rune) bool { return r == '-' || r == '_' })
	b := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		b += sandbox.Deps.Stringsdeps.ToUpper(part[:1]) + part[1:]
	}
	if b == "" {
		return "Field"
	}
	return b
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
func quietField(sandbox *api.Sandbox, conf *commandconf.CommandConf) string {
	for _, flag := range conf.Flags {
		if flag.Key == "quiet" && flag.Type == "boolean" && !flag.Array {
			return exportedName(sandbox, flag.Key)
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

func defaultLiteral(sandbox *api.Sandbox, kind, value string) string {
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
		return sandbox.Deps.Stringsdeps.Quote(value)
	}
}

func goStringList(sandbox *api.Sandbox, values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, sandbox.Deps.Stringsdeps.Quote(value))
	}
	return sandbox.Deps.Stringsdeps.Join(quoted, ", ")
}

// matchExpr builds the boolean switch guard that matches the command verb,
// e.g. `action == "version" || action == "--version"`.
func matchExpr(sandbox *api.Sandbox, identifiers []string) string {
	if len(identifiers) == 0 {
		return "false"
	}
	parts := make([]string, 0, len(identifiers))
	for _, id := range identifiers {
		parts = append(parts, "action == "+sandbox.Deps.Stringsdeps.Quote(id))
	}
	return sandbox.Deps.Stringsdeps.Join(parts, " || ")
}
