package set_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetArgInternal parses the target command's command.yaml, rebuilds the arg
// called props.Name with the changes applied — through the constructor add-arg
// calls, so an edited arg and one declared outright are the same bytes — and
// writes the file back. A rename that takes an id the command already uses is
// refused.
func SetArgInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.ArgEditProps) error {
	if utils.CommandArgEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-arg: nothing to change (pass --rename, --start, --end, --type, --required, --default, --trigger, --trigger-type, --trigger-negate, --trigger-ignore-case, --description or --clear)")
	}

	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}

	index := utils.FindCommandArg(conf, utils.CommandEntryId(sandbox, props.Name))
	if index < 0 {
		return sandbox.Deps.Std.Errorf("command %q has no arg named %q", props.Command, props.Name)
	}

	current := conf.Args[index]
	edited, err := utils.CommandArgEdited(sandbox, current, props)
	if err != nil {
		return err
	}
	if edited.Id != current.Id && utils.CommandIdTaken(conf, edited.Id) {
		return sandbox.Deps.Std.Errorf("command %q already has an arg or a flag named %s", props.Command, edited.Id)
	}

	sandbox.Deps.Std.Log("set-arg updating %s of %s \n", current.Id, utils.CommandConfPath(sandbox, utils.ResolveCommandName(sandbox, io, props.Command)))

	conf.Args[index] = edited
	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}
