package set_command

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// SetCommandInternal parses the target command's entries.yaml, overwrites
// every command-level key the caller supplied (empty strings are "leave as
// is"; --identifier and --example append) and writes the file back.
func SetCommandInternal(deps *deps.Deps, io *smartio.SmartIO, props api.CommandProps) error {
	conf, err := utils.LoadCommandConf(deps, io, props.Command)
	if err != nil {
		return err
	}
	if props.Hidden && props.Visible {
		return deps.Std.Errorf("--hidden and --visible are mutually exclusive")
	}

	changed := false
	if help := strings.TrimSpace(props.Help); help != "" {
		conf.Help, changed = help, true
	}
	if category := strings.TrimSpace(props.Category); category != "" {
		conf.Category, changed = category, true
	}
	if long := strings.TrimSpace(props.LongDescription); long != "" {
		conf.LongDescription, changed = long, true
	}
	if props.Hidden {
		conf.Hidden, changed = true, true
	}
	if props.Visible {
		conf.Hidden, changed = false, true
	}
	if len(props.Identifiers) > 0 {
		conf.Identifiers, changed = utils.AppendUnique(conf.Identifiers, props.Identifiers), true
	}
	if len(props.Examples) > 0 {
		conf.Examples, changed = utils.AppendUnique(conf.Examples, props.Examples), true
	}
	if !changed {
		return deps.Std.Errorf("set-command: nothing to change (pass --help, --category, --long-description, --hidden, --visible, --identifier or --example)")
	}

	deps.Std.Printf("set-command updating %s \n", utils.CommandEntriesPath(props.Command))

	return utils.SaveCommandConf(deps, io, props.Command, conf)
}
