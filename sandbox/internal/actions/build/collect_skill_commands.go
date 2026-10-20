package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
)

// CollectSkillCommands renders the generator's own commands — Cli.Commands of
// the binary running this build — into the pages its skill prints. They are
// read off the binary, never off the tree: a project the generator scaffolded
// carries none of the generator's command.yaml files, and the skill teaches
// what the binary that rendered it answers.
//
// Each command is turned back into the declaration it was generated from, so
// the pages are grouped, crossed with the middlewares and labelled exactly as
// docs/Commands is. Hidden commands are skipped, as they are there, and the
// rest are taken by category, then by package — the order docs/Commands walks
// its folders in, since a command's folder is its category snake_cased.
func CollectSkillCommands(sandbox *api.Sandbox) CommandDocs {
	entries := []commandDocEntry{}
	for _, command := range sandbox.Cli.Commands {
		if command.Hidden {
			continue
		}
		entries = append(entries, commandDocEntry{Name: command.Name, Conf: commandConfOf(sandbox, command)})
	}
	sandbox.Deps.SortDeps.Slice(entries, func(i, j int) bool {
		if entries[i].Conf.Category != entries[j].Conf.Category {
			return entries[i].Conf.Category < entries[j].Conf.Category
		}
		return entries[i].Name < entries[j].Name
	})
	return commandDocsOf(sandbox, entries)
}

// commandConfOf is the declaration one generated command is a 1:1 image of.
func commandConfOf(sandbox *api.Sandbox, command *api.Command) *commandconf.CommandConf {
	conf := commandconf.NewEmpty(sandbox)
	conf.Priority = command.Priority
	conf.Segments = command.Segments
	conf.HasSegments = command.Segments > 0
	conf.Strict = command.Strict
	conf.Category = command.Category
	conf.Summary = command.Summary
	conf.Description = command.Description
	conf.Examples = command.Examples
	conf.Hidden = command.Hidden

	for _, arg := range command.Args {
		conf.Args = append(conf.Args, commandconf.Arg{
			Id:          arg.Id,
			Start:       arg.Start,
			End:         arg.End,
			Type:        argTypeName(arg.Type),
			Required:    arg.Required,
			Default:     arg.Default,
			HasDefault:  arg.HasDefault,
			Description: arg.Description,
			Trigger:     triggerConfOf(arg.Trigger),
		})
	}
	for _, flag := range command.Flags {
		conf.Flags = append(conf.Flags, commandconf.Flag{
			Id:          flag.Id,
			Keys:        flag.Keys,
			HasKeys:     len(flag.Keys) > 0,
			Type:        flagTypeName(flag.Type),
			Required:    flag.Required,
			Default:     flag.Default,
			HasDefault:  flag.HasDefault,
			Min:         flag.Min,
			HasMin:      flag.HasMin,
			Max:         flag.Max,
			HasMax:      flag.HasMax,
			Enum:        flag.Enum,
			Pattern:     flag.Pattern,
			Description: flag.Description,
			Trigger:     triggerConfOf(flag.Trigger),
		})
	}
	return conf
}

// triggerConfOf is the declared trigger a generated api.Trigger spells.
func triggerConfOf(trigger api.Trigger) triggerconf.Trigger {
	return triggerconf.Trigger{
		Set:        trigger.Set,
		Type:       triggerTypeName(trigger.Type),
		Value:      trigger.Value,
		Values:     trigger.Values,
		Negate:     trigger.Negate,
		IgnoreCase: trigger.IgnoreCase,
	}
}

// argTypeName is the arg type an api.ArgType constant spells; argTypeConst
// read backwards.
func argTypeName(kind api.ArgType) string {
	switch kind {
	case api.ArgInteger:
		return "integer"
	case api.ArgNumber:
		return "number"
	case api.ArgUuid:
		return "uuid"
	}
	return commandconf.DefaultArgType
}

// flagTypeName is the flag type an api.FlagType constant spells;
// flagTypeConst read backwards.
func flagTypeName(kind api.FlagType) string {
	switch kind {
	case api.FlagInteger:
		return "integer"
	case api.FlagNumber:
		return "number"
	case api.FlagBoolean:
		return "boolean"
	case api.FlagStringArray:
		return "string-array"
	case api.FlagIntegerArray:
		return "integer-array"
	}
	return commandconf.DefaultFlagType
}

// triggerTypeName is the trigger type an api.TriggerType constant spells;
// triggerConst read backwards.
func triggerTypeName(kind api.TriggerType) string {
	switch kind {
	case api.TriggerPrefix:
		return "prefix"
	case api.TriggerTextPrefix:
		return "text-prefix"
	case api.TriggerSuffix:
		return "suffix"
	case api.TriggerRegex:
		return triggerconf.Regex
	case api.TriggerOneOf:
		return triggerconf.OneOf
	}
	return "equal"
}
