package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
)

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
	kind := sandbox.Deps.StringsDeps.ToLower(sandbox.Deps.StringsDeps.TrimSpace(raw))
	if canonical, is := TriggerAliases[kind]; is {
		kind = canonical
	}
	if !contains(triggerconf.TriggerTypes, kind) {
		return "", sandbox.Deps.StdDeps.Errorf("unknown trigger type %q (use one of %s, or starts-with, ends-with, exact, matches, any-of)",
			raw, sandbox.Deps.StringsDeps.Join(triggerconf.TriggerTypes, ", "))
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
	value := sandbox.Deps.StringsDeps.TrimSpace(props.Value)
	raw_kind := sandbox.Deps.StringsDeps.TrimSpace(props.Type)

	if value == "" {
		if raw_kind != "" {
			return triggerconf.Trigger{}, sandbox.Deps.StdDeps.Errorf("--trigger-type %q needs a --trigger to compare against", raw_kind)
		}
		if props.Negate || props.IgnoreCase {
			return triggerconf.Trigger{}, sandbox.Deps.StdDeps.Errorf("--trigger-negate and --trigger-ignore-case need a --trigger to apply to")
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
		Set:        true,
		Type:       kind,
		Values:     []string{},
		Negate:     props.Negate,
		IgnoreCase: props.IgnoreCase,
	}

	switch {
	case kind == triggerconf.Regex:
		if _, err := sandbox.Deps.StringsDeps.MatchPattern(value, ""); err != nil {
			return triggerconf.Trigger{}, sandbox.Deps.StdDeps.Errorf("invalid regex trigger %q: %s", value, err.Error())
		}
		trigger.Value = value
	case kind == triggerconf.OneOf:
		for _, raw := range sandbox.Deps.StringsDeps.Split(value, TriggerValuesSeparator) {
			one := sandbox.Deps.StringsDeps.TrimSpace(raw)
			if one == "" {
				continue
			}
			if props.OnPath {
				one = "/" + sandbox.Deps.StringsDeps.TrimLeft(one, "/")
			}
			trigger.Values = append(trigger.Values, one)
		}
	case props.OnPath && kind != "suffix":
		trigger.Value = "/" + sandbox.Deps.StringsDeps.TrimLeft(value, "/")
	default:
		trigger.Value = value
	}
	return trigger, nil
}

// MatchTrigger is the MatchTrigger of the OpinionatedAgnosCli lib
// (assets/adapter-catalog/OpinionatedAgnosCli/.../trigger.go), read against a
// declaration: a change to one is a change to the other.
func MatchTrigger(sandbox *api.Sandbox, trigger triggerconf.Trigger, text string, segmented bool) bool {
	value := trigger.Value
	if trigger.IgnoreCase && trigger.Type != triggerconf.Regex {
		value = sandbox.Deps.StringsDeps.ToLower(value)
		text = sandbox.Deps.StringsDeps.ToLower(text)
	}

	matched := text == value
	switch trigger.Type {
	case "prefix":
		if !segmented {
			matched = sandbox.Deps.StringsDeps.HasPrefix(text, value)
			break
		}
		separator := triggerSeparator(sandbox, text, value)
		value = sandbox.Deps.StringsDeps.TrimSuffix(value, separator)
		matched = value == "" || text == value || sandbox.Deps.StringsDeps.HasPrefix(text, value+separator)
	case "text-prefix":
		matched = sandbox.Deps.StringsDeps.HasPrefix(text, value)
	case "suffix":
		matched = sandbox.Deps.StringsDeps.HasSuffix(text, value)
	case triggerconf.Regex:
		if trigger.IgnoreCase {
			value = "(?i)" + value
		}
		ok, err := sandbox.Deps.StringsDeps.MatchPattern(value, text)
		matched = err == nil && ok
	case triggerconf.OneOf:
		matched = false
		for _, one := range trigger.Values {
			if trigger.IgnoreCase {
				one = sandbox.Deps.StringsDeps.ToLower(one)
			}
			if one == text {
				matched = true
				break
			}
		}
	}
	return matched != trigger.Negate
}

// triggerSeparator is segmentSeparator of the OpinionatedAgnosCli lib: "/"
// on a route's path, " " on a command's segments.
func triggerSeparator(sandbox *api.Sandbox, text string, value string) string {
	if sandbox.Deps.StringsDeps.HasPrefix(text, "/") || sandbox.Deps.StringsDeps.HasPrefix(value, "/") {
		return "/"
	}
	return " "
}

// DescribeTrigger words one trigger the way explain-route and explain-command
// print it.
func DescribeTrigger(sandbox *api.Sandbox, trigger triggerconf.Trigger) string {
	text := sandbox.Deps.StdDeps.Sprintf("%s %q", trigger.Type, trigger.Value)
	if trigger.Type == triggerconf.OneOf {
		text = sandbox.Deps.StdDeps.Sprintf("%s %q", trigger.Type, sandbox.Deps.StringsDeps.Join(trigger.Values, TriggerValuesSeparator))
	}
	if trigger.IgnoreCase {
		text += " ignoring case"
	}
	if trigger.Negate {
		text = "not " + text
	}
	return text
}
