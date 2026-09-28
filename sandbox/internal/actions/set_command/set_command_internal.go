package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetCommandInternal parses the target command's command.yaml, overwrites
// every command-level key the caller supplied (empty strings are "leave as
// is"; --identifier and --example append) and writes the file back. A further
// --identifier is one more verb the arg on segment 0 answers to: its equal
// trigger becomes a one-of.
func SetCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.SetCommandProps) error {
	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}
	if props.Hidden && props.Visible {
		return sandbox.Deps.Std.Errorf("--hidden and --visible are mutually exclusive")
	}

	changed := false
	if help := sandbox.Deps.Stringsdeps.TrimSpace(props.Help); help != "" {
		conf.Help, changed = help, true
	}
	if category := sandbox.Deps.Stringsdeps.TrimSpace(props.Category); category != "" {
		conf.Category, changed = category, true
	}
	if long := sandbox.Deps.Stringsdeps.TrimSpace(props.LongDescription); long != "" {
		conf.LongDescription, changed = long, true
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
		return sandbox.Deps.Std.Errorf("set-command: nothing to change (pass --help, --category, --long-description, --hidden, --visible, --identifier or --example)")
	}

	sandbox.Deps.Std.Log("set-command updating %s \n", utils.CommandConfPath(sandbox, utils.ResolveCommandName(sandbox, io, props.Command)))

	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}

// addVerbs adds further verbs to the arg reading segment 0 alone, which is
// what answers to the command's name.
func addVerbs(sandbox *api.Sandbox, conf *commandconf.CommandConf, verbs []string) error {
	for index, arg := range conf.Args {
		if arg.Start != 0 || arg.End != 0 || !arg.Trigger.Exists || arg.Trigger.Negate {
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
	return sandbox.Deps.Std.Errorf("--identifier needs an arg on segment 0 with an equal or one-of trigger to add the verb to")
}
