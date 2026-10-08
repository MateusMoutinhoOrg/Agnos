package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CommandClearKeys is every key --clear may take off a command.
var CommandClearKeys = []string{"segments", "examples"}

// SetCommandInternal parses the target command's command.yaml, overwrites
// every command-level key the caller supplied (empty strings are "leave as
// is"; --identifier and --example append; --before and --after land it one
// rung from another command) and writes the file back. A further
// --identifier is one more verb the arg on segment 0 answers to: its equal
// trigger becomes a one-of.
func SetCommandInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.SetCommandProps) error {
	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}
	if props.Hidden && props.Visible {
		return sandbox.Deps.StdDeps.Errorf("--hidden and --visible are mutually exclusive")
	}
	if props.Strict && props.Loose {
		return sandbox.Deps.StdDeps.Errorf("--strict and --loose are mutually exclusive")
	}

	changed := false
	if help := sandbox.Deps.StringsDeps.TrimSpace(props.Summary); help != "" {
		conf.Summary, changed = help, true
	}
	if category := sandbox.Deps.StringsDeps.TrimSpace(props.Category); category != "" {
		conf.Category, changed = category, true
	}
	if long := sandbox.Deps.StringsDeps.TrimSpace(props.Description); long != "" {
		conf.Description, changed = long, true
	}
	cleared, err := utils.RouteClearSet(sandbox, props.Clear, CommandClearKeys)
	if err != nil {
		return err
	}
	if cleared["segments"] {
		conf.Segments, conf.HasSegments, changed = 0, false, true
	}
	if cleared["examples"] {
		conf.Examples, changed = []string{}, true
	}

	relative, has_relative, err := utils.CommandRelativePriority(sandbox, io, props.Before, props.After)
	if err != nil {
		return err
	}
	if has_relative && props.HasPriority {
		return sandbox.Deps.StdDeps.Errorf("--priority excludes --before and --after: name the rung, or the command it sits next to")
	}
	if has_relative {
		props.Priority, props.HasPriority = relative, true
	}
	if props.HasPriority {
		if props.Priority < 0 {
			return sandbox.Deps.StdDeps.Errorf("--priority %d is negative: the chain runs from zero upwards", props.Priority)
		}
		conf.Priority, conf.HasPriority, changed = props.Priority, true, true
	}
	if props.HasSegments {
		if props.Segments < 1 {
			return sandbox.Deps.StdDeps.Errorf("--segments %d is below 1: clear it with --clear segments for a command that takes any count", props.Segments)
		}
		conf.Segments, conf.HasSegments, changed = props.Segments, true, true
	}
	if props.Strict {
		conf.Strict, changed = true, true
	}
	if props.Loose {
		conf.Strict, changed = false, true
	}
	if props.Hidden {
		conf.Hidden, changed = true, true
	}
	if props.Visible {
		conf.Hidden, changed = false, true
	}
	if len(props.Identifiers) > 0 {
		if err := addVerbs(sandbox, conf, props.Identifiers); err != nil {
			return err
		}
		changed = true
	}
	if len(props.Examples) > 0 {
		conf.Examples, changed = utils.AppendUnique(conf.Examples, props.Examples), true
	}
	if !changed {
		return sandbox.Deps.StdDeps.Errorf("set-command: nothing to change (pass --help, --category, --description, --priority, --before, --after, --segments, --strict, --loose, --clear, --hidden, --visible, --identifier or --example)")
	}

	sandbox.Deps.StdDeps.Logf("set-command updating %s \n", utils.CommandConfPath(sandbox, io, props.Command))

	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}

// addVerbs adds further verbs to the arg reading segment 0 alone, which is
// what answers to the command's name.
func addVerbs(sandbox *api.Sandbox, conf *commandconf.CommandConf, verbs []string) error {
	for index, arg := range conf.Args {
		if arg.Start != 0 || arg.End != 0 || !arg.Trigger.Set || arg.Trigger.Negate {
			continue
		}
		values := conf.Identifiers()
		if arg.Trigger.Type != "equal" && arg.Trigger.Type != triggerconf.OneOf {
			break
		}
		conf.Args[index].Trigger.Type = triggerconf.OneOf
		conf.Args[index].Trigger.Value = ""
		conf.Args[index].Trigger.Values = utils.AppendUnique(values, verbs)
		return nil
	}
	return sandbox.Deps.StdDeps.Errorf("--identifier needs an arg on segment 0 with an equal or one-of trigger to add the verb to")
}
