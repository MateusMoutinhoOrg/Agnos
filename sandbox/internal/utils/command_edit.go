package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
)

// The edit helpers below are what `set-arg` and `set-flag` are, the way
// route_edit.go is what `set-path` and `set-parameter` are: a declaration
// already on disk read back as the command line that would have written it,
// the keys being changed written over that, and the whole built again by the
// same constructor the add- side calls.

// CommandArgClearKeys is every key --clear may take off an arg.
var CommandArgClearKeys = []string{"trigger", "trigger-negate", "trigger-ignore-case", "type", "required", "default", "description"}

// CommandFlagClearKeys is every key --clear may take off a flag.
var CommandFlagClearKeys = []string{"keys", "type", "required", "default", "min", "max", "enum", "pattern", "trigger", "trigger-negate", "trigger-ignore-case", "description"}

// triggerTyped is one declared trigger as the command line types it: a
// one-of's values joined by the separator NewTrigger splits them on.
func triggerTyped(sandbox *api.Sandbox, trigger triggerconf.Trigger) string {
	if trigger.Type == triggerconf.OneOf {
		return sandbox.Deps.StringsDeps.Join(trigger.Values, TriggerValuesSeparator)
	}
	return trigger.Value
}

// EditCommandArg rebuilds one arg with the changes applied, holding the
// result to every rule NewCommandArg holds a new one to.
func EditCommandArg(sandbox *api.Sandbox, current commandconf.Arg, props api.SetArgProps) (commandconf.Arg, error) {
	strs := sandbox.Deps.StringsDeps
	cleared, err := RouteClearSet(sandbox, props.Clear, CommandArgClearKeys)
	if err != nil {
		return commandconf.Arg{}, err
	}

	built := api.AddArgProps{
		Name:        current.Id,
		Start:       strs.FormatInt(int64(current.Start), 10),
		End:         strs.FormatInt(int64(current.End), 10),
		Type:        current.Type,
		Required:    current.Required,
		Description: current.Description,
	}
	if current.HasDefault {
		built.Default = current.Default
	}
	if current.Trigger.Set {
		built.TriggerType, built.Trigger = current.Trigger.Type, triggerTyped(sandbox, current.Trigger)
		built.TriggerNegate, built.TriggerIgnoreCase = current.Trigger.Negate, current.Trigger.IgnoreCase
	}

	if cleared["trigger"] {
		built.TriggerType, built.Trigger = "", ""
		built.TriggerNegate, built.TriggerIgnoreCase = false, false
	}
	if cleared["trigger-negate"] {
		built.TriggerNegate = false
	}
	if cleared["trigger-ignore-case"] {
		built.TriggerIgnoreCase = false
	}
	if cleared["type"] {
		built.Type = ""
	}
	if cleared["required"] {
		built.Required = false
	}
	if cleared["default"] {
		built.Default = ""
	}
	if cleared["description"] {
		built.Description = ""
	}

	if value := strs.TrimSpace(props.Rename); value != "" {
		built.Name = value
	}
	if value := strs.TrimSpace(props.Start); value != "" {
		built.Start = value
	}
	if value := strs.TrimSpace(props.End); value != "" {
		built.End = value
	}
	if value := strs.TrimSpace(props.Type); value != "" {
		built.Type = value
	}
	if value := strs.TrimSpace(props.Default); value != "" {
		built.Default, built.Required = value, false
	}
	if props.Required {
		built.Required, built.Default = true, ""
	}
	if value := strs.TrimSpace(props.Description); value != "" {
		built.Description = value
	}
	if value := strs.TrimSpace(props.Trigger); value != "" {
		built.Trigger = value
	}
	if value := strs.TrimSpace(props.TriggerType); value != "" {
		built.TriggerType = value
	}
	if props.TriggerNegate {
		built.TriggerNegate = true
	}
	if props.TriggerIgnoreCase {
		built.TriggerIgnoreCase = true
	}

	return NewCommandArg(sandbox, built, current.Start)
}

// IsCommandArgEditEmpty reports an edit that changes nothing.
func IsCommandArgEditEmpty(sandbox *api.Sandbox, props api.SetArgProps) bool {
	given := sandbox.Deps.StringsDeps.TrimSpace(props.Rename + props.Start + props.End + props.Type +
		props.Default + props.TriggerType + props.Trigger + props.Description)
	return given == "" && len(props.Clear) == 0 && !props.Required && !props.TriggerNegate && !props.TriggerIgnoreCase
}

// EditCommandFlag rebuilds one flag with the changes applied, holding the
// result to every rule NewCommandFlag holds a new one to.
func EditCommandFlag(sandbox *api.Sandbox, current commandconf.Flag, props api.SetFlagProps) (commandconf.Flag, error) {
	strs := sandbox.Deps.StringsDeps
	cleared, err := RouteClearSet(sandbox, props.Clear, CommandFlagClearKeys)
	if err != nil {
		return commandconf.Flag{}, err
	}

	built := api.AddFlagProps{
		Name:        current.Id,
		Type:        current.Type,
		Required:    current.Required,
		Enum:        current.Enum,
		Pattern:     current.Pattern,
		Description: current.Description,
	}
	if current.HasKeys {
		built.Keys = current.Keys
	}
	if current.HasDefault {
		built.Default = current.Default
	}
	if current.HasMin {
		built.Min = strs.FormatFloat(current.Min, 'g', -1, 64)
	}
	if current.HasMax {
		built.Max = strs.FormatFloat(current.Max, 'g', -1, 64)
	}
	if current.Trigger.Set {
		built.TriggerType, built.Trigger = current.Trigger.Type, triggerTyped(sandbox, current.Trigger)
		built.TriggerNegate, built.TriggerIgnoreCase = current.Trigger.Negate, current.Trigger.IgnoreCase
	}

	for key := range cleared {
		switch key {
		case "keys":
			built.Keys = nil
		case "type":
			built.Type = ""
		case "required":
			built.Required = false
		case "default":
			built.Default = ""
		case "min":
			built.Min = ""
		case "max":
			built.Max = ""
		case "enum":
			built.Enum = nil
		case "pattern":
			built.Pattern = ""
		case "trigger":
			built.TriggerType, built.Trigger = "", ""
			built.TriggerNegate, built.TriggerIgnoreCase = false, false
		case "trigger-negate":
			built.TriggerNegate = false
		case "trigger-ignore-case":
			built.TriggerIgnoreCase = false
		case "description":
			built.Description = ""
		}
	}

	if value := strs.TrimSpace(props.Rename); value != "" {
		built.Name = value
		if !current.HasKeys && len(props.Keys) == 0 {
			built.Keys = nil
		}
	}
	if len(props.Keys) > 0 {
		built.Keys = props.Keys
	}
	if value := strs.TrimSpace(props.Type); value != "" {
		built.Type = value
	}
	if value := strs.TrimSpace(props.Default); value != "" {
		built.Default, built.Required = value, false
	}
	if props.Required {
		built.Required, built.Default = true, ""
	}
	if value := strs.TrimSpace(props.Min); value != "" {
		built.Min = value
	}
	if value := strs.TrimSpace(props.Max); value != "" {
		built.Max = value
	}
	if len(props.Enum) > 0 {
		built.Enum = props.Enum
	}
	if value := strs.TrimSpace(props.Pattern); value != "" {
		built.Pattern = value
	}
	if value := strs.TrimSpace(props.Description); value != "" {
		built.Description = value
	}
	if value := strs.TrimSpace(props.Trigger); value != "" {
		built.Trigger = value
	}
	if value := strs.TrimSpace(props.TriggerType); value != "" {
		built.TriggerType = value
	}
	if props.TriggerNegate {
		built.TriggerNegate = true
	}
	if props.TriggerIgnoreCase {
		built.TriggerIgnoreCase = true
	}

	return NewCommandFlag(sandbox, built)
}

// IsCommandFlagEditEmpty reports an edit that changes nothing.
func IsCommandFlagEditEmpty(sandbox *api.Sandbox, props api.SetFlagProps) bool {
	given := sandbox.Deps.StringsDeps.TrimSpace(props.Rename + props.Type + props.Default + props.Min + props.Max +
		props.Pattern + props.TriggerType + props.Trigger + props.Description)
	return given == "" && len(props.Keys) == 0 && len(props.Enum) == 0 && len(props.Clear) == 0 &&
		!props.Required && !props.TriggerNegate && !props.TriggerIgnoreCase
}
