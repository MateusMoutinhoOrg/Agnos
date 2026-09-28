package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *CommandConf) {
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
	conf.Pattern = func() string {
		return Pattern(sandbox, conf)
	}
	conf.Identifiers = func() []string {
		return Identifiers(sandbox, conf)
	}
}

// Identifiers are the literal texts the arg reading from segment 0 answers to:
// the value of an equal trigger, or every value of a one-of. A command whose
// first arg is no literal — a middleware, a capture — has none.
func Identifiers(sandbox *api.Sandbox, conf *CommandConf) []string {
	for _, arg := range conf.Args {
		if arg.Start != 0 || !arg.Trigger.Exists || arg.Trigger.Negate {
			continue
		}
		switch arg.Trigger.Type {
		case "equal":
			return []string{arg.Trigger.Value}
		case "one-of":
			return append([]string{}, arg.Trigger.Values...)
		}
	}
	return []string{}
}

// Pattern draws the command line a command answers, one slot per segment,
// filled by the first arg that reaches it:
//
//	equal    its value                 add-flag
//	one-of   (a|b)                     (help|h)
//	prefix   its value, then …         route …
//	capture  <Id>, <Id:type>, or <Id…> to the end
//	other    text-prefix Value*, suffix *Value, regex ~(Value)
//
// A negated trigger is appended as !(…) rather than drawn. A command reading
// nothing at all — a middleware on every line — reads as "*".
func Pattern(sandbox *api.Sandbox, conf *CommandConf) string {
	slots := map[int]string{}
	last := -1
	tail := ""
	tailAt := -1
	negated := ""

	place := func(index int, text string) {
		if _, taken := slots[index]; taken {
			return
		}
		slots[index] = text
		if index > last {
			last = index
		}
	}
	placeTail := func(index int, text string) {
		if tailAt < 0 || index < tailAt {
			tail, tailAt = text, index
		}
	}

	for _, arg := range conf.Args {
		trigger := arg.Trigger
		if trigger.Exists && trigger.Negate {
			negated += " !(" + triggerText(sandbox, trigger) + ")"
			continue
		}

		if !trigger.Exists {
			label := arg.Id
			if arg.Type != "" && arg.Type != DefaultArgType {
				label += ":" + arg.Type
			}
			if arg.End == LastSegment || arg.End > arg.Start {
				label += "…"
			}
			if !arg.Required {
				label = "[" + label + "]"
			} else {
				label = "<" + label + ">"
			}
			if arg.End == LastSegment {
				placeTail(arg.Start, label)
				continue
			}
			place(arg.Start, label)
			continue
		}

		switch trigger.Type {
		case "equal", "prefix":
			index := arg.Start
			for _, word := range sandbox.Deps.Stringsdeps.Fields(trigger.Value) {
				place(index, word)
				index++
			}
			if trigger.Type == "prefix" {
				placeTail(index, "…")
			}
		case "one-of":
			place(arg.Start, "("+triggerText(sandbox, trigger)+")")
		case "text-prefix":
			placeTail(arg.Start, trigger.Value+"*")
		case "suffix":
			placeTail(arg.Start, "*"+trigger.Value)
		case "regex":
			placeTail(arg.Start, "~("+trigger.Value+")")
		}
	}

	pieces := []string{}
	tailed := false
	for index := 0; index <= last; index++ {
		if text, ok := slots[index]; ok {
			pieces = append(pieces, text)
			continue
		}
		if tailAt >= 0 && index >= tailAt {
			pieces, tailed = append(pieces, tail), true
			break
		}
		pieces = append(pieces, "*")
	}
	if tailAt >= 0 && !tailed {
		for len(pieces) < tailAt {
			pieces = append(pieces, "*")
		}
		pieces = append(pieces, tail)
	}

	if len(pieces) == 0 || (len(pieces) == 1 && pieces[0] == "…") {
		return "*" + negated
	}
	return sandbox.Deps.Stringsdeps.Join(pieces, " ") + negated
}

// triggerText is what a trigger compares against, as a pattern draws it: its
// value, or a one-of's values joined by "|".
func triggerText(sandbox *api.Sandbox, trigger Trigger) string {
	if trigger.Type != "one-of" {
		return trigger.Value
	}
	return sandbox.Deps.Stringsdeps.Join(trigger.Values, "|")
}
