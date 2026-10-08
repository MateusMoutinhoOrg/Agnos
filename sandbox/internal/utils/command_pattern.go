package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
)

// CommandPattern is what one --pattern compiles to: the entries of `args` it
// stands for, and the segment count a command line has to have — HasSegments
// false when the pattern ends on a {*rest}, which takes any count from there
// on.
type CommandPattern struct {
	Args        []commandconf.Arg
	Segments    int
	HasSegments bool
}

// CompileCommandPattern turns one command-line shape typed on the command line
// into the args command.yaml declares. The yaml never holds the pattern
// itself: it holds what the pattern means, and conf.Pattern() draws it back.
//
//	route add                  {Command 0..1 equal "route add"}
//	route add {name}           + {Name 2..2}
//	scale {count:integer}      + {Count 1..1 type integer}
//	exec {*rest}               + {Rest 1..-1}
//
// Literal words in a row are one arg, the first of them named Command and the
// rest after their words. Without a {*…} the pattern fixes the segment count.
func CompileCommandPattern(sandbox *api.Sandbox, raw string) (CommandPattern, error) {
	words := sandbox.Deps.StringsDeps.Fields(raw)
	if len(words) == 0 {
		return CommandPattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q declares no word", raw)
	}

	compiled := CommandPattern{Args: []commandconf.Arg{}, Segments: len(words), HasSegments: true}
	taken := map[string]bool{}
	for _, reserved := range CommandReservedIds {
		taken[reserved] = true
	}
	claim := func(id string, index int) (string, error) {
		if id == "" {
			return "", sandbox.Deps.StdDeps.Errorf("--pattern %q names an empty capture at segment %d", raw, index)
		}
		if taken[id] {
			return "", sandbox.Deps.StdDeps.Errorf("--pattern %q binds %s twice, or one Input already carries", raw, id)
		}
		taken[id] = true
		return id, nil
	}

	literal_start := -1
	literal := []string{}
	flush := func(end int) error {
		if len(literal) == 0 {
			return nil
		}
		id := "Command"
		if literal_start > 0 || taken[id] {
			id = GoIdentifier(sandbox, sandbox.Deps.StringsDeps.Join(literal, "-"))
		}
		if id == "" || id[0] < 'A' || id[0] > 'Z' || taken[id] {
			id = "Seg" + sandbox.Deps.StringsDeps.FormatInt(int64(literal_start), 10)
		}
		id, err := claim(id, literal_start)
		if err != nil {
			return err
		}
		compiled.Args = append(compiled.Args, commandconf.Arg{
			Id:      id,
			Start:   literal_start,
			End:     end,
			Type:    commandconf.DefaultArgType,
			Trigger: triggerconf.Trigger{Set: true, Type: "equal", Value: sandbox.Deps.StringsDeps.Join(literal, " "), Values: []string{}},
		})
		literal, literal_start = []string{}, -1
		return nil
	}

	for index, word := range words {
		if !sandbox.Deps.StringsDeps.HasPrefix(word, "{") {
			if sandbox.Deps.StringsDeps.Contains(word, "{") || sandbox.Deps.StringsDeps.Contains(word, "}") {
				return CommandPattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q mixes text and a capture in the word %q: a capture is a whole word", raw, word)
			}
			if sandbox.Deps.StringsDeps.HasPrefix(word, "-") {
				return CommandPattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q holds %q: a segment never starts with -, declare it with add-flag", raw, word)
			}
			if literal_start < 0 {
				literal_start = index
			}
			literal = append(literal, word)
			continue
		}

		if err := flush(index - 1); err != nil {
			return CommandPattern{}, err
		}
		if !sandbox.Deps.StringsDeps.HasSuffix(word, "}") {
			return CommandPattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q opens a capture it does not close: %q", raw, word)
		}
		inner := word[1 : len(word)-1]

		if sandbox.Deps.StringsDeps.HasPrefix(inner, "*") {
			if index != len(words)-1 {
				return CommandPattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q puts %s before the end: a {*…} capture takes the rest of the line", raw, word)
			}
			id, err := claim(GoIdentifier(sandbox, inner[1:]), index)
			if err != nil {
				return CommandPattern{}, err
			}
			compiled.Args = append(compiled.Args, commandconf.Arg{
				Id: id, Start: index, End: commandconf.LastSegment, Type: commandconf.DefaultArgType,
			})
			compiled.Segments, compiled.HasSegments = 0, false
			continue
		}

		name, kind := inner, commandconf.DefaultArgType
		if parts := sandbox.Deps.StringsDeps.Split(inner, ":"); len(parts) == 2 {
			name, kind = parts[0], sandbox.Deps.StringsDeps.ToLower(parts[1])
		}
		if !contains(commandconf.ArgTypes, kind) {
			return CommandPattern{}, sandbox.Deps.StdDeps.Errorf("--pattern %q declares the unknown type %q (use one of %s)",
				raw, kind, sandbox.Deps.StringsDeps.Join(commandconf.ArgTypes, ", "))
		}
		id, err := claim(GoIdentifier(sandbox, name), index)
		if err != nil {
			return CommandPattern{}, err
		}
		compiled.Args = append(compiled.Args, commandconf.Arg{Id: id, Start: index, End: index, Type: kind, Required: true})
	}

	if err := flush(len(words) - 1); err != nil {
		return CommandPattern{}, err
	}
	return compiled, nil
}
