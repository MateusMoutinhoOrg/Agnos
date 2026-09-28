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

// CommandMatch is what one command makes of one command line, read the way
// the generated IsActionable and CommandHandler read it: whether it runs,
// why not when it does not, and — for one that runs — the usage error its
// values raise before its handler does, "" when they bind.
type CommandMatch struct {
	Runs    bool
	Reason  string
	Failure string
}

// commandEndOfFlags is endOfFlags of the generated IsActionable.go.
const commandEndOfFlags = "--"

// SplitCommandArgv is SplitArgv of the generated IsActionable.go: the
// segments of a command line and the index in argv of each.
func SplitCommandArgv(sandbox *api.Sandbox, argv []string) ([]string, []int) {
	segments := []string{}
	indices := []int{}

	index := 0
	for ; index < len(argv); index++ {
		if sandbox.Deps.Stringsdeps.HasPrefix(argv[index], "-") {
			break
		}
		segments = append(segments, argv[index])
		indices = append(indices, index)
	}
	for ; index < len(argv); index++ {
		if argv[index] != commandEndOfFlags {
			continue
		}
		for rest := index + 1; rest < len(argv); rest++ {
			segments = append(segments, argv[rest])
			indices = append(indices, rest)
		}
		break
	}
	return segments, indices
}

// commandFlagsEnd is FlagsEnd of the generated IsActionable.go.
func commandFlagsEnd(argv []string) int {
	for index, token := range argv {
		if token == commandEndOfFlags {
			return index
		}
	}
	return len(argv)
}

// commandArgSlice is ArgSlice of the generated IsActionable.go.
func commandArgSlice(segments []string, arg commandconf.Arg) ([]string, bool) {
	end := arg.End
	if end < 0 {
		if arg.Start >= len(segments) {
			return []string{}, arg.Start == len(segments)
		}
		end = len(segments) - 1
	}
	if arg.Start < 0 || arg.Start > end || end >= len(segments) {
		return nil, false
	}
	return segments[arg.Start : end+1], true
}

// commandArgConverts is ArgValue of the generated IsActionable.go, reduced to
// whether the slice converts.
func commandArgConverts(sandbox *api.Sandbox, arg commandconf.Arg, values []string) bool {
	if arg.End != arg.Start {
		return true
	}
	if len(values) != 1 {
		return false
	}
	return CheckCommandLiteral(sandbox, arg.Type, "", values[0]) == nil
}

// commandFlagOccurrences is where one flag's keys stand on a command line, up
// to its bare "--": the index of each key, which is the index of its value
// minus one for every flag but a boolean.
func commandFlagOccurrences(argv []string, flag commandconf.Flag) []int {
	found := []int{}
	for index, token := range argv[:commandFlagsEnd(argv)] {
		if contains(flag.Keys, token) {
			found = append(found, index)
		}
	}
	return found
}

// MatchCommandArgv is the generated IsActionable and CommandHandler of one
// command read against its command.yaml: whether it runs for argv and, when
// it does, the usage error its values raise. The tokens it reads are marked
// on consumed — its segments only when it is strict, as the dispatch does —
// and a strict command reports the first token nobody of the chain read.
func MatchCommandArgv(sandbox *api.Sandbox, conf *commandconf.CommandConf, argv []string, consumed []bool) CommandMatch {
	strs := sandbox.Deps.Stringsdeps
	segments, indices := SplitCommandArgv(sandbox, argv)

	if conf.HasSegments && len(segments) != conf.Segments {
		return CommandMatch{Reason: sandbox.Deps.Std.Sprintf("the line has %d segments, the command takes %d", len(segments), conf.Segments)}
	}
	for _, arg := range conf.Args {
		values, found := commandArgSlice(segments, arg)
		if !found {
			if arg.Trigger.Exists {
				return CommandMatch{Reason: sandbox.Deps.Std.Sprintf("arg %s: the line has no segment %d", arg.Id, arg.Start)}
			}
			continue
		}
		text := strs.Join(values, " ")
		if !commandArgConverts(sandbox, arg, values) {
			return CommandMatch{Reason: sandbox.Deps.Std.Sprintf("arg %s: %q is not a %s", arg.Id, text, arg.Type)}
		}
		if arg.Trigger.Exists && !MatchTrigger(sandbox, arg.Trigger, text, true) {
			return CommandMatch{Reason: sandbox.Deps.Std.Sprintf("arg %s: %q fails %s", arg.Id, text, DescribeTrigger(sandbox, arg.Trigger))}
		}
	}
	for _, flag := range conf.Flags {
		if !flag.Trigger.Exists {
			continue
		}
		value := ""
		occurrences := commandFlagOccurrences(argv, flag)
		switch {
		case len(occurrences) == 0:
			return CommandMatch{Reason: sandbox.Deps.Std.Sprintf("flag %s: not on the line, and it declares %s", flag.Keys[0], DescribeTrigger(sandbox, flag.Trigger))}
		case flag.Type == "boolean":
			value = "true"
		case occurrences[0]+1 < len(argv):
			value = argv[occurrences[0]+1]
		}
		if !MatchTrigger(sandbox, flag.Trigger, value, false) {
			return CommandMatch{Reason: sandbox.Deps.Std.Sprintf("flag %s: %q fails %s", flag.Keys[0], value, DescribeTrigger(sandbox, flag.Trigger))}
		}
	}

	match := CommandMatch{Runs: true}

	for _, arg := range conf.Args {
		values, found := commandArgSlice(segments, arg)
		if !found || len(values) == 0 {
			if arg.Required && match.Failure == "" {
				match.Failure = sandbox.Deps.Std.Sprintf("a usage error: required arg '%s' not provided", arg.Id)
			}
			continue
		}
		if !conf.Strict {
			continue
		}
		end := arg.End
		if end < 0 {
			end = len(segments) - 1
		}
		for index := arg.Start; index <= end; index++ {
			consumed[indices[index]] = true
		}
	}
	for _, flag := range conf.Flags {
		occurrences := commandFlagOccurrences(argv, flag)
		if len(occurrences) == 0 && flag.Required && match.Failure == "" {
			match.Failure = sandbox.Deps.Std.Sprintf("a usage error: required flag '%s' not provided", flag.Keys[0])
		}
		for _, index := range occurrences {
			consumed[index] = true
			if flag.Type == "boolean" {
				continue
			}
			if index+1 >= len(argv) {
				if match.Failure == "" {
					match.Failure = sandbox.Deps.Std.Sprintf("a usage error: flag '%s' expects a value", flag.Id)
				}
				continue
			}
			consumed[index+1] = true
			if problem := commandFlagProblem(sandbox, flag, argv[index+1]); problem != "" && match.Failure == "" {
				match.Failure = "a usage error: flag '" + flag.Keys[0] + "': " + problem
			}
		}
	}
	if end := commandFlagsEnd(argv); end < len(argv) {
		consumed[end] = true
	}

	if conf.Strict && match.Failure == "" {
		for index, token := range argv {
			if consumed[index] {
				continue
			}
			if strs.HasPrefix(token, "-") {
				match.Failure = sandbox.Deps.Std.Sprintf("a usage error: unknown flag %q", token)
			} else {
				match.Failure = sandbox.Deps.Std.Sprintf("a usage error: unexpected argument %q", token)
			}
			break
		}
	}
	return match
}

// commandFlagProblem is flagValue of the generated CommandHandler.go read
// against a declaration: what is wrong with one raw value of a flag — its
// enum, its pattern, its type, its bounds — "" when it binds.
func commandFlagProblem(sandbox *api.Sandbox, flag commandconf.Flag, raw string) string {
	if len(flag.Enum) > 0 && !contains(flag.Enum, raw) {
		return sandbox.Deps.Std.Sprintf("%q is not one of %s", raw, sandbox.Deps.Stringsdeps.Join(flag.Enum, ", "))
	}
	if flag.Pattern != "" {
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(flag.Pattern, raw)
		if err != nil || !matched {
			return sandbox.Deps.Std.Sprintf("%q does not match %s", raw, flag.Pattern)
		}
	}
	if err := CheckCommandLiteral(sandbox, flag.Type, "the value", raw); err != nil {
		return err.Error()
	}
	if flag.Type != "integer" && flag.Type != "number" && flag.Type != "integer-array" {
		return ""
	}
	number, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
	if err != nil {
		return ""
	}
	if flag.HasMin && number < flag.Min {
		return "must be >= " + sandbox.Deps.Stringsdeps.FormatFloat(flag.Min, 'g', -1, 64)
	}
	if flag.HasMax && number > flag.Max {
		return "must be <= " + sandbox.Deps.Stringsdeps.FormatFloat(flag.Max, 'g', -1, 64)
	}
	return ""
}
