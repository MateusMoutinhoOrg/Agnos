package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectCommands reads every command.yaml under sandbox/internal/commands — a
// directory holding one is a command, at any depth — and returns one data map per command, in run order — lowest `priority` first,
// then by name — for the generated.new.go and generated.input.go of each command and
// the {{range .Commands}} loop of sandbox/internal/cli/generated.new.go.
// What the map holds is the declaration itself: the dispatch and the help
// screens read it back off Cli.Commands at runtime. `help` is collected like
// every other command: its command.yaml is written by GenerateHelpCommandYaml
// just before this runs. It is the server's CollectRoutes, for the cli.
func CollectCommands(sandbox *api.Sandbox, io *stagedfs.StagedFS) ([]map[string]any, error) {
	var commands []map[string]any

	units := utils.CommandDirs(sandbox, io)
	if err := utils.CheckUniqueUnitNames(sandbox, "command", units); err != nil {
		return nil, err
	}

	for _, unit := range units {
		content, err := io.ReadFile(unit.Dir + "/" + utils.CommandConfFile)
		if err != nil {
			continue
		}

		conf, err := commandconf.New(sandbox, string(content))
		if err != nil {
			return nil, sandbox.Deps.StdDeps.Errorf("%s/%s: %w", unit.Dir, utils.CommandConfFile, err)
		}

		commands = append(commands, commandData(sandbox, unit, conf))
	}

	sortCommands(sandbox, commands)
	return commands, nil
}

// sortCommands puts the commands in the order the dispatch runs them: the
// lowest `priority` first, then by name for a stable tie-break. The ordering is
// the collector's, so Cli.Commands is already in run order and the dispatch
// only has to range.
func sortCommands(sandbox *api.Sandbox, commands []map[string]any) {
	sandbox.Deps.SortDeps.SliceStable(commands, func(i int, j int) bool {
		left, right := commands[i], commands[j]
		if left["Priority"] != right["Priority"] {
			return left["Priority"].(int) < right["Priority"].(int)
		}
		return left["CommandName"].(string) < right["CommandName"].(string)
	})
}

// commandData is one command as the generated.new.go and generated.input.go read it.
// Dir is the project-relative directory the command sits in, which the
// generated cli imports it from.
func commandData(sandbox *api.Sandbox, unit utils.UnitDir, conf *commandconf.CommandConf) map[string]any {
	args := make([]map[string]any, 0, len(conf.Args))
	for _, arg := range conf.Args {
		args = append(args, argData(sandbox, arg))
	}
	flags := make([]map[string]any, 0, len(conf.Flags))
	for _, flag := range conf.Flags {
		flags = append(flags, flagData(sandbox, flag))
	}

	segments := 0
	if conf.HasSegments {
		segments = conf.Segments
	}

	return map[string]any{
		"CommandName": unit.Name,
		"Dir":         unit.Dir,
		"Identifiers": conf.Identifiers(),
		"Priority":    conf.Priority,
		"Segments":    segments,
		"Strict":      conf.Strict,
		"Pattern":     conf.Pattern(),
		"Category":    conf.Category,
		"Summary":     conf.Summary,
		"Description": conf.Description,
		"Examples":    conf.Examples,
		"Hidden":      conf.Hidden,
		"Args":        args,
		"Flags":       flags,
	}
}

// argData is one entry of `args` as the generated api.CommandArg literal and
// the Input field read it.
func argData(sandbox *api.Sandbox, arg commandconf.Arg) map[string]any {
	return map[string]any{
		"Id":          arg.Id,
		"Start":       arg.Start,
		"End":         arg.End,
		"Type":        argTypeConst(arg.Type),
		"GoType":      argGoType(arg),
		"Required":    arg.Required,
		"Default":     arg.Default,
		"HasDefault":  arg.HasDefault,
		"Description": arg.Description,
		"Trigger":     triggerData(arg.Trigger),
	}
}

// flagData is one entry of `flags` as the generated api.CommandFlag literal
// and the Input field read it.
func flagData(sandbox *api.Sandbox, flag commandconf.Flag) map[string]any {
	return map[string]any{
		"Id":          flag.Id,
		"Keys":        flag.Keys,
		"Type":        flagTypeConst(flag.Type),
		"GoType":      flagGoType(flag.Type),
		"Required":    flag.Required,
		"Default":     flag.Default,
		"HasDefault":  flag.HasDefault,
		"Min":         numberLabel(sandbox, flag.Type, flag.Min, flag.HasMin),
		"Max":         numberLabel(sandbox, flag.Type, flag.Max, flag.HasMax),
		"HasMin":      flag.HasMin,
		"HasMax":      flag.HasMax,
		"Enum":        flag.Enum,
		"Pattern":     flag.Pattern,
		"Description": flag.Description,
		"Trigger":     triggerData(flag.Trigger),
	}
}

// argTypeConst is the api.ArgType constant an arg type is spelled with.
func argTypeConst(kind string) string {
	switch kind {
	case "integer":
		return "api.ArgInteger"
	case "number":
		return "api.ArgNumber"
	case "uuid":
		return "api.ArgUuid"
	}
	return "api.ArgString"
}

// argGoType is the Go type an arg binds as: one segment in its type, a range
// as a []string.
func argGoType(arg commandconf.Arg) string {
	if arg.End != arg.Start {
		return "[]string"
	}
	switch arg.Type {
	case "integer":
		return "int"
	case "number":
		return "float64"
	}
	return "string"
}

// flagTypeConst is the api.FlagType constant a flag type is spelled with.
func flagTypeConst(kind string) string {
	switch kind {
	case "integer":
		return "api.FlagInteger"
	case "number":
		return "api.FlagNumber"
	case "boolean":
		return "api.FlagBoolean"
	case "string-array":
		return "api.FlagStringArray"
	case "integer-array":
		return "api.FlagIntegerArray"
	}
	return "api.FlagString"
}

// flagGoType is the Go type a flag binds as.
func flagGoType(kind string) string {
	switch kind {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "string-array":
		return "[]string"
	case "integer-array":
		return "[]int"
	}
	return "string"
}

// numberLabel renders a min/max bound as the literal it has in command.yaml
// ("" when the bound is unset).
func numberLabel(sandbox *api.Sandbox, kind string, value float64, has bool) string {
	if !has {
		return ""
	}
	if kind == "integer" || kind == "integer-array" {
		return sandbox.Deps.StringsDeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.StringsDeps.FormatFloat(value, 'g', -1, 64)
}

// ─── helpers ────────────────────────────────────────────────────────────────

func lastSegmentOf(sandbox *api.Sandbox, path string) string {
	parts := sandbox.Deps.StringsDeps.Split(path, "/")
	return parts[len(parts)-1]
}

func exportedName(sandbox *api.Sandbox, raw string) string {
	if name := utils.GoIdentifier(sandbox, raw); name != "" {
		return name
	}
	return "Field"
}

func goStringList(sandbox *api.Sandbox, values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, sandbox.Deps.StringsDeps.Quote(value))
	}
	return sandbox.Deps.StringsDeps.Join(quoted, ", ")
}
