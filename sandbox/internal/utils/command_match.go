package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
)

// Reach is whether a middleware runs in front of a command, read off the two
// declarations alone — without a command line.
type Reach int

const (
	// NoReach is a middleware that never runs in front of the command.
	NoReach Reach = iota
	// Runs is a middleware that runs in front of every line of the command.
	Runs
	// MayRun is a middleware whose trigger reads a segment the command
	// captures, or a regex: whether it runs depends on what is typed, and
	// explain-command gives the exact answer.
	MayRun
)

// MiddlewareReach crosses one middleware with one command: whether it runs in
// front of it — a lower rung, and every trigger of its args holding on the
// command's literal words — and, when one of its flags declares a trigger, the
// condition that adds ("only when --x equal \"y\"").
func MiddlewareReach(sandbox *api.Sandbox, middleware *commandconf.CommandConf, command *commandconf.CommandConf) (Reach, string) {
	if middleware.Strict || middleware.Priority >= command.Priority {
		return NoReach, ""
	}

	reach := Runs
	if middleware.HasSegments && (!command.HasSegments || command.Segments != middleware.Segments) {
		if command.HasSegments {
			return NoReach, ""
		}
		reach = MayRun
	}

	words := literalWords(sandbox, command)
	for _, arg := range middleware.Args {
		if !arg.Trigger.Exists {
			continue
		}
		switch triggerReach(sandbox, arg, words) {
		case NoReach:
			return NoReach, ""
		case MayRun:
			reach = MayRun
		}
	}

	conditions := []string{}
	for _, flag := range middleware.Flags {
		if flag.Trigger.Exists {
			conditions = append(conditions, flag.Keys[0]+" "+DescribeTrigger(sandbox, flag.Trigger))
		}
	}
	if len(conditions) == 0 {
		return reach, ""
	}
	return reach, "only when " + sandbox.Deps.Stringsdeps.Join(conditions, " and ")
}

// literalWords is, segment by segment, the one word a command's line is known
// to carry there: what an equal trigger of its args spells. A segment a
// capture, a one-of or anything else reads is left out — it is not known.
func literalWords(sandbox *api.Sandbox, command *commandconf.CommandConf) map[int]string {
	words := map[int]string{}
	for _, arg := range command.Args {
		if !arg.Trigger.Exists || arg.Trigger.Negate || arg.Trigger.Type != "equal" {
			continue
		}
		for offset, word := range sandbox.Deps.Stringsdeps.Fields(arg.Trigger.Value) {
			if _, known := words[arg.Start+offset]; !known {
				words[arg.Start+offset] = word
			}
		}
	}
	return words
}

// triggerReach is whether one trigger of a middleware's arg holds on a
// command's literal words: decided when every word it compares is known, or
// when the known ones already settle a prefix; MayRun otherwise.
func triggerReach(sandbox *api.Sandbox, arg commandconf.Arg, words map[int]string) Reach {
	trigger := arg.Trigger
	if trigger.Type == triggerconf.Regex {
		return MayRun
	}

	known := []string{}
	complete := arg.End != commandconf.LastSegment
	for index := arg.Start; arg.End == commandconf.LastSegment || index <= arg.End; index++ {
		word, has := words[index]
		if !has {
			if arg.End != commandconf.LastSegment {
				complete = false
			}
			break
		}
		known = append(known, word)
	}
	if arg.End == commandconf.LastSegment {
		// The command's line may go on past its literal words, so the text
		// is only known when the trigger settles on them.
		complete = false
	}

	text := sandbox.Deps.Stringsdeps.Join(known, " ")
	if complete {
		return boolReach(MatchTrigger(sandbox, trigger, text, true))
	}

	if trigger.Type == "prefix" && !trigger.Negate {
		value := sandbox.Deps.Stringsdeps.Fields(trigger.Value)
		if len(value) <= len(known) {
			return boolReach(MatchTrigger(sandbox, trigger, text, true))
		}
		for index, word := range known {
			if value[index] != word {
				return NoReach
			}
		}
	}
	return MayRun
}

// boolReach is Runs for a trigger that holds, NoReach for one that fails.
func boolReach(holds bool) Reach {
	if holds {
		return Runs
	}
	return NoReach
}
