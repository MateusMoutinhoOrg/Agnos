package routeconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, conf *RouteConf) {
	conf.Render = func() string {
		return Render(sandbox, conf)
	}
	conf.Pattern = func() string {
		return Pattern(sandbox, conf)
	}
	conf.SchemaJson = func() string {
		return SchemaJson(sandbox, conf)
	}
}

// Pattern is the route's path as docs and messages spell it, one slot per
// request segment, filled by the first entry of `paths` that reaches it:
//
//	equal   its value's segments                /get-article
//	prefix  its value's segments, then /*       /admin/*
//	capture {Id}, {Id:type}, or {*Id} to the end  /get-article/{article:integer}
//	other   text-prefix Value*, suffix *Value, regex ~(Value), one-of (a|b)
//
// A slot nothing reaches reads as *, and a negated trigger is appended as
// !(…) rather than drawn, since it names what the path is not. A route with no
// path at all reads as "/".
func Pattern(sandbox *api.Sandbox, conf *RouteConf) string {
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

	for _, path := range conf.Paths {
		if path.Trigger.Exists && path.Trigger.Negate {
			negated += " !(" + triggerText(sandbox, path.Trigger) + ")"
			continue
		}

		if !path.Trigger.Exists {
			label := path.Id
			if path.Type != "" && path.Type != DefaultPathType {
				label += ":" + path.Type
			}
			if path.End == LastSegment {
				placeTail(path.Start, "{*"+path.Id+"}")
				continue
			}
			for index := path.Start; index <= path.End; index++ {
				place(index, "{"+label+"}")
			}
			continue
		}

		switch path.Trigger.Type {
		case "equal", "prefix":
			index := path.Start
			for _, segment := range sandbox.Deps.Stringsdeps.Split(path.Trigger.Value, "/") {
				if segment == "" {
					continue
				}
				place(index, segment)
				index++
			}
			if path.Trigger.Type == "prefix" {
				placeTail(index, "*")
			}
		case "text-prefix":
			placeTail(path.Start, sandbox.Deps.Stringsdeps.TrimLeft(path.Trigger.Value, "/")+"*")
		case "suffix":
			placeTail(path.Start, "*"+path.Trigger.Value)
		case "regex":
			placeTail(path.Start, "~("+path.Trigger.Value+")")
		case "one-of":
			place(path.Start, "("+triggerText(sandbox, path.Trigger)+")")
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

	pattern := "/" + sandbox.Deps.Stringsdeps.Join(pieces, "/")
	return pattern + negated
}

// triggerText is what a trigger compares against, as a pattern draws it: its
// value, or a one-of's values joined by "|", each without its leading "/".
func triggerText(sandbox *api.Sandbox, trigger Trigger) string {
	if trigger.Type != "one-of" {
		return trigger.Value
	}
	values := []string{}
	for _, value := range trigger.Values {
		values = append(values, sandbox.Deps.Stringsdeps.TrimLeft(value, "/"))
	}
	return sandbox.Deps.Stringsdeps.Join(values, "|")
}
