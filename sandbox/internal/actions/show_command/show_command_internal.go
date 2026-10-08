package show_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// The indent the tree is drawn with, as show-route draws it.
const branch = "  "

// ShowCommandInternal reads one command.yaml and renders it as the lines of a
// tree: the command line it answers, what it says about itself, where it sits
// in the chain, its args, its flags, and the middlewares that run in front of
// it with the flags they add. It writes nothing.
func ShowCommandInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, command string) ([]string, error) {
	name := utils.ResolveCommandName(sandbox, io, command)
	conf, err := utils.LoadCommandConf(sandbox, io, name)
	if err != nil {
		return nil, err
	}
	chain, err := utils.LoadCommandChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	lines := []string{conf.Pattern()}
	if conf.Summary != "" {
		lines = append(lines, branch+conf.Summary)
	}
	if conf.Category != "" {
		lines = append(lines, sandbox.Deps.StdDeps.Sprintf("%scategory  %s", branch, conf.Category))
	}
	lines = append(lines, sandbox.Deps.StdDeps.Sprintf("%spriority  %d", branch, conf.Priority))
	if conf.HasSegments {
		lines = append(lines, sandbox.Deps.StdDeps.Sprintf("%ssegments  %d", branch, conf.Segments))
	}
	if !conf.Strict {
		lines = append(lines, branch+"middleware, not strict")
	}
	if conf.Hidden {
		lines = append(lines, branch+"hidden")
	}
	lines = append(lines, chainLine(sandbox, chain, name)...)

	if len(conf.Args) > 0 {
		lines = append(lines, "", "args")
		for _, arg := range conf.Args {
			lines = append(lines, branch+argLine(sandbox, arg))
		}
	}
	if len(conf.Flags) > 0 {
		lines = append(lines, "", "flags")
		for _, flag := range conf.Flags {
			lines = append(lines, branch+flagLine(sandbox, flag))
		}
	}

	if conf.Strict {
		front := []string{}
		for _, entry := range chain {
			reach, condition := utils.CommandMiddlewareReach(sandbox, entry.Conf, conf)
			if reach == utils.NoReach {
				continue
			}
			text := sandbox.Deps.StdDeps.Sprintf("%-20s priority %d", utils.CommandName(sandbox, entry.Name), entry.Conf.Priority)
			notes := []string{}
			if reach == utils.MayRun {
				notes = append(notes, "may run")
			}
			if condition != "" {
				notes = append(notes, condition)
			}
			front = append(front, branch+withNotes(sandbox, text, notes))
			for _, flag := range entry.Conf.Flags {
				if utils.CommandKeyTaken(conf, flag.Keys) == "" {
					front = append(front, branch+branch+flagLine(sandbox, flag))
				}
			}
		}
		if len(front) > 0 {
			lines = append(lines, "", "middlewares in front")
			lines = append(lines, front...)
		}
	}

	return lines, nil
}

// chainLine is where the command sits in the chain: its place in run order,
// and the commands on either side of it.
func chainLine(sandbox *api.Sandbox, chain []utils.CommandChainEntry, name string) []string {
	for index, entry := range chain {
		if entry.Name != name {
			continue
		}
		text := sandbox.Deps.StdDeps.Sprintf("%schain     #%d of %d", branch, index+1, len(chain))
		if index > 0 {
			text += ", after " + utils.CommandName(sandbox, chain[index-1].Name)
		}
		if index < len(chain)-1 {
			text += ", before " + utils.CommandName(sandbox, chain[index+1].Name)
		}
		return []string{text}
	}
	return []string{}
}

// argLine is one arg: the Input field it binds, the segments it reads, and
// every rule on them.
func argLine(sandbox *api.Sandbox, arg commandconf.Arg) string {
	text := sandbox.Deps.StdDeps.Sprintf("%-20s segments %d..%d", arg.Id, arg.Start, arg.End)
	notes := []string{}
	if arg.Type != "" && arg.Type != commandconf.DefaultArgType {
		notes = append(notes, arg.Type)
	}
	if arg.Trigger.Set {
		notes = append(notes, utils.DescribeTrigger(sandbox, arg.Trigger))
	}
	if arg.Required {
		notes = append(notes, "required")
	}
	if arg.HasDefault {
		notes = append(notes, sandbox.Deps.StdDeps.Sprintf("default %q", arg.Default))
	}
	if arg.Description != "" {
		notes = append(notes, arg.Description)
	}
	return withNotes(sandbox, text, notes)
}

// flagLine is one flag: the Input field it binds, the keys it is typed
// under, and every rule on its value.
func flagLine(sandbox *api.Sandbox, flag commandconf.Flag) string {
	text := sandbox.Deps.StdDeps.Sprintf("%-20s %s", flag.Id, sandbox.Deps.StringsDeps.Join(flag.Keys, ", "))
	notes := []string{flag.Type}
	if flag.Required {
		notes = append(notes, "required")
	}
	if flag.HasDefault {
		notes = append(notes, sandbox.Deps.StdDeps.Sprintf("default %q", flag.Default))
	}
	if flag.HasMin {
		notes = append(notes, "min "+sandbox.Deps.StringsDeps.FormatFloat(flag.Min, 'g', -1, 64))
	}
	if flag.HasMax {
		notes = append(notes, "max "+sandbox.Deps.StringsDeps.FormatFloat(flag.Max, 'g', -1, 64))
	}
	if len(flag.Enum) > 0 {
		notes = append(notes, "one of "+sandbox.Deps.StringsDeps.Join(flag.Enum, "/"))
	}
	if flag.Pattern != "" {
		notes = append(notes, "pattern "+flag.Pattern)
	}
	if flag.Trigger.Set {
		notes = append(notes, utils.DescribeTrigger(sandbox, flag.Trigger))
	}
	if flag.Description != "" {
		notes = append(notes, flag.Description)
	}
	return withNotes(sandbox, text, notes)
}

// withNotes appends the notes of one line after its text, bracketed.
func withNotes(sandbox *api.Sandbox, text string, notes []string) string {
	if len(notes) == 0 {
		return text
	}
	return text + "  (" + sandbox.Deps.StringsDeps.Join(notes, "; ") + ")"
}
