package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
)

// The edit helpers below are what `set-path`, `set-parameter` and
// `set-body-field` are: a declaration already on disk read back as the command
// line that would have written it, the keys being changed written over that,
// and the whole built again by the same constructor the add- side calls. A
// declaration that was edited and one that was declared outright are the same
// bytes, because they went through the same function.
//
// An empty value means "leave it as it is", so taking a key off again is what
// --clear is for: one repeated flag naming the keys to drop, applied before
// the changes so that clearing and setting the same key in one line reads as
// the change.

// RoutePathClearKeys is every key --clear may take off an entry of `paths`.
var RoutePathClearKeys = []string{"trigger", "trigger-negate", "trigger-ignore-case", "type", "description"}

// RouteParameterClearKeys is every key --clear may take off an entry of
// `parameters`.
var RouteParameterClearKeys = []string{"description", "examples", "default", "required", "trigger", "trigger-negate", "trigger-ignore-case"}

// RouteBodyFieldClearKeys is every keyword --clear may take off a body
// property.
var RouteBodyFieldClearKeys = []string{
	"required", "array", "min", "max", "exclusive-min", "exclusive-max",
	"format", "pattern", "enum", "const", "nullable",
	"min-items", "max-items", "unique-items", "additional-properties",
}

// RouteClearSet reads a repeated --clear into the set of keys it names,
// refusing one the command carries nothing to clear for — a typo there would
// otherwise read as "nothing to do" and leave the declaration as it was.
func RouteClearSet(sandbox *api.Sandbox, clear []string, known []string) (map[string]bool, error) {
	cleared := map[string]bool{}

	for _, raw := range clear {
		key := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
		if key == "" {
			continue
		}

		found := false
		for _, one := range known {
			if one == key {
				found = true
				break
			}
		}
		if !found {
			return nil, sandbox.Deps.Std.Errorf("--clear %q is not a key that can be taken off (use %s)", raw, sandbox.Deps.Stringsdeps.Join(known, ", "))
		}

		cleared[key] = true
	}

	return cleared, nil
}

// RoutePathEdited rebuilds one entry of `paths` with the changes applied,
// holding the result to every rule NewRoutePath holds a new one to.
func RoutePathEdited(sandbox *api.Sandbox, current routeconf.Path, props api.RoutePathEditProps) (routeconf.Path, error) {
	cleared, err := RouteClearSet(sandbox, props.Clear, RoutePathClearKeys)
	if err != nil {
		return routeconf.Path{}, err
	}

	built := api.RoutePathProps{
		Id:          current.Id,
		Start:       sandbox.Deps.Stringsdeps.FormatInt(int64(current.Start), 10),
		End:         sandbox.Deps.Stringsdeps.FormatInt(int64(current.End), 10),
		Type:        current.Type,
		Description: current.Description,
	}
	if current.Trigger.Exists {
		built.TriggerType, built.Trigger = current.Trigger.Type, current.Trigger.Value
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
	if cleared["description"] {
		built.Description = ""
	}

	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Rename); value != "" {
		built.Id = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Start); value != "" {
		built.Start = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.End); value != "" {
		built.End = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Description); value != "" {
		built.Description = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Trigger); value != "" {
		built.Trigger = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.TriggerType); value != "" {
		built.TriggerType = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Type); value != "" {
		built.Type = value
	}
	if props.TriggerNegate {
		built.TriggerNegate = true
	}
	if props.TriggerIgnoreCase {
		built.TriggerIgnoreCase = true
	}

	return NewRoutePath(sandbox, built)
}

// RoutePathEditEmpty reports an edit that changes nothing, so the command can
// say so instead of rewriting a file with the bytes already in it.
func RoutePathEditEmpty(sandbox *api.Sandbox, props api.RoutePathEditProps) bool {
	given := sandbox.Deps.Stringsdeps.TrimSpace(props.Rename + props.Start + props.End + props.Type +
		props.TriggerType + props.Trigger + props.Description)
	return given == "" && len(props.Clear) == 0 && !props.TriggerNegate && !props.TriggerIgnoreCase
}

// RouteParameterEdited rebuilds one entry of `parameters` with the changes
// applied, holding the result to every rule NewRouteParameter holds a new one
// to.
func RouteParameterEdited(sandbox *api.Sandbox, current routeconf.Parameter, props api.RouteParameterEditProps) (routeconf.Parameter, error) {
	cleared, err := RouteClearSet(sandbox, props.Clear, RouteParameterClearKeys)
	if err != nil {
		return routeconf.Parameter{}, err
	}

	built := api.RouteParameterProps{
		Name:        current.Key,
		Type:        current.Type,
		Fonts:       current.Fonts,
		Required:    current.Required,
		Description: current.Description,
		Examples:    current.Examples,
	}
	if current.HasDefault {
		built.Default = current.Default
	}
	if current.Trigger.Exists {
		built.TriggerType, built.Trigger = current.Trigger.Type, current.Trigger.Value
		built.TriggerNegate, built.TriggerIgnoreCase = current.Trigger.Negate, current.Trigger.IgnoreCase
	}

	if cleared["description"] {
		built.Description = ""
	}
	if cleared["trigger-negate"] {
		built.TriggerNegate = false
	}
	if cleared["trigger-ignore-case"] {
		built.TriggerIgnoreCase = false
	}
	if cleared["examples"] {
		built.Examples = nil
	}
	if cleared["default"] {
		built.Default = ""
	}
	if cleared["required"] {
		built.Required = false
	}
	if cleared["trigger"] {
		built.TriggerType, built.Trigger = "", ""
		built.TriggerNegate, built.TriggerIgnoreCase = false, false
	}

	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Rename); value != "" {
		built.Name = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Type); value != "" {
		built.Type = value
	}
	if len(props.Fonts) > 0 {
		built.Fonts = props.Fonts
	}
	if props.Required {
		built.Required = true
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Default); value != "" {
		built.Default = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Description); value != "" {
		built.Description = value
	}
	if len(props.Examples) > 0 {
		built.Examples = AppendUnique(built.Examples, props.Examples)
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Trigger); value != "" {
		built.Trigger = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.TriggerType); value != "" {
		built.TriggerType = value
	}
	if props.TriggerNegate {
		built.TriggerNegate = true
	}
	if props.TriggerIgnoreCase {
		built.TriggerIgnoreCase = true
	}

	return NewRouteParameter(sandbox, built)
}

// RouteParameterEditEmpty reports an edit that changes nothing.
func RouteParameterEditEmpty(sandbox *api.Sandbox, props api.RouteParameterEditProps) bool {
	given := sandbox.Deps.Stringsdeps.TrimSpace(props.Rename + props.Type + props.Default +
		props.TriggerType + props.Trigger + props.Description)
	return given == "" && len(props.Fonts) == 0 && len(props.Examples) == 0 && len(props.Clear) == 0 &&
		!props.Required && !props.TriggerNegate && !props.TriggerIgnoreCase
}

// RouteBodyFieldEdit is one edited body property: the schema it becomes,
// whether its parent object still demands it, and the keywords a change of
// --type left behind — a format on a property that is no longer text has
// nowhere to go, and dropping it silently would be the one thing an edit must
// not do.
type RouteBodyFieldEdit struct {
	Schema   *routeconf.Schema
	Required bool
	Dropped  []string
}

// RouteBodyFieldEdited rebuilds one declared body property with the changes
// applied, through the same RouteBodyPropertySchema that declares a new one.
func RouteBodyFieldEdited(sandbox *api.Sandbox, current *routeconf.Schema, required bool, props api.RouteBodyFieldEditProps) (RouteBodyFieldEdit, error) {
	cleared, err := RouteClearSet(sandbox, props.Clear, RouteBodyFieldClearKeys)
	if err != nil {
		return RouteBodyFieldEdit{}, err
	}

	built := RouteBodyFieldPropsOf(sandbox, current)
	built.Required = required

	if cleared["required"] {
		built.Required = false
	}
	if cleared["array"] {
		built.Array, built.MinItems, built.MaxItems, built.UniqueItems = false, "", "", false
	}
	if cleared["min"] {
		built.Min = ""
	}
	if cleared["max"] {
		built.Max = ""
	}
	if cleared["exclusive-min"] {
		built.ExclusiveMin = ""
	}
	if cleared["exclusive-max"] {
		built.ExclusiveMax = ""
	}
	if cleared["format"] {
		built.Format = ""
	}
	if cleared["pattern"] {
		built.Pattern = ""
	}
	if cleared["enum"] {
		built.Enum = nil
	}
	if cleared["const"] {
		built.Const = ""
	}
	if cleared["nullable"] {
		built.Nullable = false
	}
	if cleared["min-items"] {
		built.MinItems = ""
	}
	if cleared["max-items"] {
		built.MaxItems = ""
	}
	if cleared["unique-items"] {
		built.UniqueItems = false
	}
	if cleared["additional-properties"] {
		built.AdditionalProperties, built.NoAdditionalProperties = false, false
	}

	retyped := ""
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Type); value != "" {
		kind, err := RouteSchemaKind(sandbox, value)
		if err != nil {
			return RouteBodyFieldEdit{}, err
		}
		if kind != built.Type {
			retyped = kind
		}
		built.Type = kind
	}

	if props.Required {
		built.Required = true
	}
	if props.Array {
		built.Array = true
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Min); value != "" {
		built.Min = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Max); value != "" {
		built.Max = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.ExclusiveMin); value != "" {
		built.ExclusiveMin = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.ExclusiveMax); value != "" {
		built.ExclusiveMax = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Format); value != "" {
		built.Format = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Pattern); value != "" {
		built.Pattern = value
	}
	if len(props.Enum) > 0 {
		built.Enum = props.Enum
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.Const); value != "" {
		built.Const = value
	}
	if props.Nullable {
		built.Nullable = true
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.MinItems); value != "" {
		built.MinItems = value
	}
	if value := sandbox.Deps.Stringsdeps.TrimSpace(props.MaxItems); value != "" {
		built.MaxItems = value
	}
	if props.UniqueItems {
		built.UniqueItems = true
	}
	if props.AdditionalProperties {
		built.AdditionalProperties, built.NoAdditionalProperties = true, false
	}
	if props.NoAdditionalProperties {
		built.NoAdditionalProperties, built.AdditionalProperties = true, false
	}

	dropped := dropRetypedKeywords(&built, retyped)

	schema, err := RouteBodyPropertySchema(sandbox, built)
	if err != nil {
		return RouteBodyFieldEdit{}, err
	}
	if err := carrySchemaChildren(sandbox, current, schema); err != nil {
		return RouteBodyFieldEdit{}, err
	}

	return RouteBodyFieldEdit{Schema: schema, Required: built.Required, Dropped: dropped}, nil
}

// carrySchemaChildren keeps what was declared inside an object property when
// that property is rewritten. RouteBodyPropertySchema builds one node out of
// keywords and knows nothing of the properties under it, so a --max-items on a
// list of objects, or a --no-additional-properties on one, would otherwise drop
// every key beneath it without saying so.
//
// A property that stops being an object is refused rather than emptied: the
// keys under it are declarations someone wrote, and losing them is not
// something an edit that was asked for a type may do on its own.
func carrySchemaChildren(sandbox *api.Sandbox, current *routeconf.Schema, rebuilt *routeconf.Schema) error {
	was := SchemaObjectOf(current)
	if was == nil || was.Type != "object" || len(was.Properties) == 0 {
		return nil
	}

	now := SchemaObjectOf(rebuilt)
	if now == nil || now.Type != "object" {
		return sandbox.Deps.Std.Errorf(
			"%d properties are declared under this one, so it cannot stop being an object: drop them first, or drop the whole of it with remove-body-field",
			len(was.Properties))
	}

	now.Properties = was.Properties
	now.Required = was.Required
	return nil
}

// dropRetypedKeywords takes off the keywords the property's new type cannot
// carry, and names every one it took — a --type that turned text into a number
// leaves a format behind, and the alternative to dropping it is an error about
// a keyword the person never typed.
func dropRetypedKeywords(built *api.RouteBodyFieldProps, retyped string) []string {
	if retyped == "" {
		return nil
	}

	dropped := []string{}
	numeric := retyped == "int" || retyped == "float"

	if retyped != "string" {
		if built.Format != "" {
			built.Format, dropped = "", append(dropped, "format")
		}
		if built.Pattern != "" {
			built.Pattern, dropped = "", append(dropped, "pattern")
		}
	}
	if !numeric {
		if built.ExclusiveMin != "" {
			built.ExclusiveMin, dropped = "", append(dropped, "exclusive-min")
		}
		if built.ExclusiveMax != "" {
			built.ExclusiveMax, dropped = "", append(dropped, "exclusive-max")
		}
	}
	if !numeric && retyped != "string" {
		if built.Min != "" {
			built.Min, dropped = "", append(dropped, "min")
		}
		if built.Max != "" {
			built.Max, dropped = "", append(dropped, "max")
		}
	}
	if retyped != "object" && (built.AdditionalProperties || built.NoAdditionalProperties) {
		built.AdditionalProperties, built.NoAdditionalProperties = false, false
		dropped = append(dropped, "additional-properties")
	}

	return dropped
}

// RouteBodyFieldEditEmpty reports an edit that changes nothing.
func RouteBodyFieldEditEmpty(sandbox *api.Sandbox, props api.RouteBodyFieldEditProps) bool {
	given := sandbox.Deps.Stringsdeps.TrimSpace(props.Rename + props.Type + props.Min + props.Max +
		props.ExclusiveMin + props.ExclusiveMax + props.Format + props.Pattern + props.Const +
		props.MinItems + props.MaxItems)
	return given == "" && len(props.Enum) == 0 && len(props.Clear) == 0 &&
		!props.Required && !props.Array && !props.Nullable && !props.UniqueItems &&
		!props.AdditionalProperties && !props.NoAdditionalProperties
}
