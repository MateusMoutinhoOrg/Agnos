package set_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetFlagInternal parses the target command's command.yaml, rebuilds the flag
// called props.Name — matched by its name, its id or one of its keys — with
// the changes applied, through the constructor add-flag calls, and writes the
// file back. A rename or new keys that another flag of the command already
// takes are refused.
func SetFlagInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.FlagEditProps) error {
	if utils.CommandFlagEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-flag: nothing to change (pass --rename, --key, --type, --required, --default, --min, --max, --enum, --pattern, --trigger, --trigger-type, --trigger-negate, --trigger-ignore-case, --description or --clear)")
	}

	conf, err := utils.LoadCommandConf(sandbox, io, props.Command)
	if err != nil {
		return err
	}

	index := utils.FindCommandFlagNamed(sandbox, conf, props.Name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("command %q has no flag named %q", props.Command, props.Name)
	}

	current := conf.Flags[index]
	edited, err := utils.CommandFlagEdited(sandbox, current, props)
	if err != nil {
		return err
	}

	others := *conf
	others.Flags = append(append([]commandconf.Flag{}, conf.Flags[:index]...), conf.Flags[index+1:]...)
	if edited.Id != current.Id && utils.CommandIdTaken(&others, edited.Id) {
		return sandbox.Deps.Std.Errorf("command %q already has an arg or a flag named %s", props.Command, edited.Id)
	}
	if taken := utils.CommandKeyTaken(&others, edited.Keys); taken != "" {
		return sandbox.Deps.Std.Errorf("key %q is already used by another flag of %q", taken, props.Command)
	}

	sandbox.Deps.Std.Log("set-flag updating %s of %s \n", current.Id, utils.CommandConfPath(sandbox, utils.ResolveCommandName(sandbox, io, props.Command)))

	conf.Flags[index] = edited
	return utils.SaveCommandConf(sandbox, io, props.Command, conf)
}
