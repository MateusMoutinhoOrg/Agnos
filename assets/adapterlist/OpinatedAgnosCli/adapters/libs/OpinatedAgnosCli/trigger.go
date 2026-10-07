package opinatedagnoscli

import (
	"regexp"
	"strings"

	opinatedagnoscli "{{.Module}}/sandbox/deps/OpinatedAgnosCli"
)

// matchTrigger reports whether one text meets a trigger: the comparison its
// Type names, run without regard to case when it declares IgnoreCase, and
// inverted when it declares Negate. Segmented is true for a text a prefix
// reads segment by segment — a route's path slice, a command's segments.
//
// It is mirrored by agnos's own sandbox/internal/utils/trigger.go: a change to
// one is a change to the other.
func matchTrigger(trigger opinatedagnoscli.Trigger, text string, segmented bool) bool {
	return compareTrigger(trigger, text, segmented) != trigger.Negate
}

// compareTrigger is the comparison of matchTrigger before Negate. On a
// segmented text a prefix holds on a segment boundary alone — "/admin" is
// "/admin" or "/admin/…", never "/administrator"; "add" is "add" or "add …",
// never "add-flag" — which is what tells it from a text-prefix. An empty prefix
// holds on every text. A value has no segments, so there the two are one.
func compareTrigger(trigger opinatedagnoscli.Trigger, text string, segmented bool) bool {
	value := trigger.Value
	if trigger.IgnoreCase && trigger.Type != opinatedagnoscli.RegexTrigger {
		value = strings.ToLower(value)
		text = strings.ToLower(text)
	}

	switch trigger.Type {
	case opinatedagnoscli.PrefixTrigger:
		if !segmented {
			return strings.HasPrefix(text, value)
		}
		separator := segmentSeparator(text, value)
		value = strings.TrimSuffix(value, separator)
		return value == "" || text == value || strings.HasPrefix(text, value+separator)
	case opinatedagnoscli.TextPrefixTrigger:
		return strings.HasPrefix(text, value)
	case opinatedagnoscli.SuffixTrigger:
		return strings.HasSuffix(text, value)
	case opinatedagnoscli.RegexTrigger:
		if trigger.IgnoreCase {
			value = "(?i)" + value
		}
		matched, err := regexp.MatchString(value, text)
		return err == nil && matched
	case opinatedagnoscli.OneOfTrigger:
		for _, one := range trigger.Values {
			if trigger.IgnoreCase {
				one = strings.ToLower(one)
			}
			if one == text {
				return true
			}
		}
		return false
	}
	return text == value
}

// segmentSeparator is what splits the segments of a segmented text: "/" on a
// route's path, which always starts with one, and " " on a command's segments,
// which never does.
func segmentSeparator(text string, value string) string {
	if strings.HasPrefix(text, "/") || strings.HasPrefix(value, "/") {
		return "/"
	}
	return " "
}
