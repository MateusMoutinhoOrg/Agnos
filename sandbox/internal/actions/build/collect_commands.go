package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectCommands reads every sandbox/internal/commands/<name>/entries.yaml and
// returns one data map per command, in listing order, for the generated
// sandbox/internal/commands/<name>/new.go and the {{range .Commands}} loop of
// sandbox/binds/cli.go. What the map holds is the declaration itself: the
// dispatch and the help screens read it back off sandbox.Commands at runtime,
// so nothing here is a Go spelling of anything. `help` is collected like every
// other command: its entries.yaml is written by GenerateHelpEntriesYaml just
// before this runs.
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
		"Identifiers":     conf.Identifiers,
		"IdentifiersGo":   goStringList(sandbox, conf.Identifiers),
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
		"Key":           field.Key,
		"Type":          field.Type,
		"IsArray":       field.Array,
		"Identifiers":   field.Identifiers,
		"IdentifiersGo": goStringList(sandbox, field.Identifiers),
		"Required":      field.Required,
		"HasDefault":    field.HasDefault,
		"Default":       field.Default,
		"Description":   field.Description,
		"Examples":      field.Examples,
		"Min":           numberLabel(sandbox, field, field.Min, field.HasMin),
		"Max":           numberLabel(sandbox, field, field.Max, field.HasMax),
		"HasMin":        field.HasMin,
		"HasMax":        field.HasMax,
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

func goStringList(sandbox *api.Sandbox, values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, sandbox.Deps.Stringsdeps.Quote(value))
	}
	return sandbox.Deps.Stringsdeps.Join(quoted, ", ")
}
