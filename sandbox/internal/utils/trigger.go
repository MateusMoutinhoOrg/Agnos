package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// TriggerDir is the generated package holding the MatchTrigger the cli and
// the server layers share: the build writes it while either one is on, and
// the purge of the last of them removes it.
const TriggerDir = GeneratedDir + "/trigger"

// RemoveTriggerUnlessUsed drops TriggerDir once the layer being purged was the
// last one matching on a trigger — other is the extension of the one left.
func RemoveTriggerUnlessUsed(sandbox *api.Sandbox, io *smartio.SmartIO, other string) error {
	used, err := ExtensionEnabled(sandbox, io, other)
	if err != nil || used || !io.IsDir(TriggerDir) {
		return err
	}
	for _, entry := range io.ListAllRecursively(TriggerDir) {
		io.RemoveDir(entry)
	}
	io.RemoveDir(TriggerDir)
	return nil
}

// TriggerAliases maps the spellings a trigger type may be typed in onto the
// one a declaration carries; the canonical names map onto themselves.
var TriggerAliases = map[string]string{
	"starts-with": "prefix",
	"ends-with":   "suffix",
	"exact":       "equal",
	"equals":      "equal",
	"matches":     "regex",
	"any-of":      triggerconf.OneOf,
}

// TriggerValuesSeparator splits the one value a one-of trigger is typed with
// on the command line into the values it accepts: "version,--version".
const TriggerValuesSeparator = ","

// NormalizeTriggerType normalizes a trigger type typed on the command line: an
// alias becomes the type it stands for, and anything that is neither is
// refused with the list of both.
func NormalizeTriggerType(sandbox *api.Sandbox, raw string) (string, error) {
	kind := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if canonical, is := TriggerAliases[kind]; is {
		kind = canonical
	}
	if !contains(triggerconf.TriggerTypes, kind) {
		return "", sandbox.Deps.Std.Errorf("unknown trigger type %q (use one of %s, or starts-with, ends-with, exact, matches, any-of)",
			raw, sandbox.Deps.Stringsdeps.Join(triggerconf.TriggerTypes, ", "))
	}
	return kind, nil
}

// TriggerProps is one trigger as it is typed on the command line: the value,
// how it is compared ("" is equal), and the two switches on it. OnPath reports
// a trigger compared against a route's path slice, which reads with the
// leading slash every slice carries — "users" and "/users" are the same
// request.
type TriggerProps struct {
	Type       string
	Value      string
	Negate     bool
	IgnoreCase bool
	OnPath     bool
}

// NewTrigger normalizes a trigger typed on the command line: its type
// defaults to equal and takes the aliases of TriggerAliases, and every type
// but suffix and regex compared against a path slice gets its leading slash. A
// regex is taken verbatim, and has to compile; a one-of splits its value on
// TriggerValuesSeparator.
func NewTrigger(sandbox *api.Sandbox, props TriggerProps) (triggerconf.Trigger, error) {
	value := sandbox.Deps.Stringsdeps.TrimSpace(props.Value)
	raw_kind := sandbox.Deps.Stringsdeps.TrimSpace(props.Type)

	if value == "" {
		if raw_kind != "" {
			return triggerconf.Trigger{}, sandbox.Deps.Std.Errorf("--trigger-type %q needs a --trigger to compare against", raw_kind)
		}
		if props.Negate || props.IgnoreCase {
			return triggerconf.Trigger{}, sandbox.Deps.Std.Errorf("--trigger-negate and --trigger-ignore-case need a --trigger to apply to")
		}
		return triggerconf.Trigger{}, nil
	}
	if raw_kind == "" {
		raw_kind = "equal"
	}
	kind, err := NormalizeTriggerType(sandbox, raw_kind)
	if err != nil {
		return triggerconf.Trigger{}, err
	}

	trigger := triggerconf.Trigger{
		Exists:     true,
		Type:       kind,
		Values:     []string{},
		Negate:     props.Negate,
		IgnoreCase: props.IgnoreCase,
	}

	switch {
	case kind == triggerconf.Regex:
		if _, err := sandbox.Deps.Stringsdeps.MatchPattern(value, ""); err != nil {
			return triggerconf.Trigger{}, sandbox.Deps.Std.Errorf("invalid regex trigger %q: %s", value, err.Error())
		}
		trigger.Value = value
	case kind == triggerconf.OneOf:
		for _, raw := range sandbox.Deps.Stringsdeps.Split(value, TriggerValuesSeparator) {
			one := sandbox.Deps.Stringsdeps.TrimSpace(raw)
			if one == "" {
				continue
			}
			if props.OnPath {
				one = "/" + sandbox.Deps.Stringsdeps.TrimLeft(one, "/")
			}
			trigger.Values = append(trigger.Values, one)
		}
	case props.OnPath && kind != "suffix":
		trigger.Value = "/" + sandbox.Deps.Stringsdeps.TrimLeft(value, "/")
	default:
		trigger.Value = value
	}
	return trigger, nil
}

// MatchTrigger is MatchTrigger of the generated
// sandbox/internal/generated/trigger/MatchTrigger.go, read against a
// declaration: a change to one is a change to the other.
func MatchTrigger(sandbox *api.Sandbox, trigger triggerconf.Trigger, text string, segmented bool) bool {
	value := trigger.Value
	if trigger.IgnoreCase && trigger.Type != triggerconf.Regex {
		value = sandbox.Deps.Stringsdeps.ToLower(value)
		text = sandbox.Deps.Stringsdeps.ToLower(text)
	}

	matched := text == value
	switch trigger.Type {
	case "prefix":
		if !segmented {
			matched = sandbox.Deps.Stringsdeps.HasPrefix(text, value)
			break
		}
		separator := triggerSeparator(sandbox, text, value)
		value = sandbox.Deps.Stringsdeps.TrimSuffix(value, separator)
		matched = value == "" || text == value || sandbox.Deps.Stringsdeps.HasPrefix(text, value+separator)
	case "text-prefix":
		matched = sandbox.Deps.Stringsdeps.HasPrefix(text, value)
	case "suffix":
		matched = sandbox.Deps.Stringsdeps.HasSuffix(text, value)
	case triggerconf.Regex:
		if trigger.IgnoreCase {
			value = "(?i)" + value
		}
		ok, err := sandbox.Deps.Stringsdeps.MatchPattern(value, text)
		matched = err == nil && ok
	case triggerconf.OneOf:
		matched = false
		for _, one := range trigger.Values {
			if trigger.IgnoreCase {
				one = sandbox.Deps.Stringsdeps.ToLower(one)
			}
			if one == text {
				matched = true
				break
			}
		}
	}
	return matched != trigger.Negate
}

// triggerSeparator is segmentSeparator of the generated MatchTrigger.go: "/"
// on a route's path, " " on a command's segments.
func triggerSeparator(sandbox *api.Sandbox, text string, value string) string {
	if sandbox.Deps.Stringsdeps.HasPrefix(text, "/") || sandbox.Deps.Stringsdeps.HasPrefix(value, "/") {
		return "/"
	}
	return " "
}

// DescribeTrigger words one trigger the way explain-route and explain-command
// print it.
func DescribeTrigger(sandbox *api.Sandbox, trigger triggerconf.Trigger) string {
	text := sandbox.Deps.Std.Sprintf("%s %q", trigger.Type, trigger.Value)
	if trigger.Type == triggerconf.OneOf {
		text = sandbox.Deps.Std.Sprintf("%s %q", trigger.Type, sandbox.Deps.Stringsdeps.Join(trigger.Values, TriggerValuesSeparator))
	}
	if trigger.IgnoreCase {
		text += " ignoring case"
	}
	if trigger.Negate {
		text = "not " + text
	}
	return text
}
