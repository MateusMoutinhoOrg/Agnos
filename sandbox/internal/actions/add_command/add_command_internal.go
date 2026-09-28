package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// defaultCategory is the heading a command lands under when it names none,
// and defaultMiddlewareCategory the one of a --middleware.
const defaultCategory = "Commands"
const defaultMiddlewareCategory = "Middleware"

// AddCommandInternal writes the two hand-written files of a new command
// package: a command.yaml and the InternalPureHandler.go stub. Its args are
// one arg answering to the command's name on segment 0 — every command line,
// for a --middleware — or another --trigger, or what a --pattern compiles to.
// It refuses to overwrite an existing command (via io.WriteFile).
func AddCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddCommandProps) error {
	strs := sandbox.Deps.Stringsdeps
	help := strs.TrimSpace(props.Help)
	if help == "" {
		return sandbox.Deps.Std.Errorf("add-command requires --help")
	}

	name := props.Name
	if err := utils.ValidateCommandName(sandbox, name); err != nil {
		return err
	}
	utils.NoteNormalizedCommandName(sandbox, name)

	identifier := utils.CommandIdentifier(sandbox, name)
	pkg := utils.CommandPackage(sandbox, name)

	if utils.IsGeneratedCommand(sandbox, name) {
		return sandbox.Deps.Std.Errorf("the %s command is generated and cannot be declared", identifier)
	}

	conf := commandconf.NewEmpty(sandbox)
	conf.Help = help

	if strs.TrimSpace(props.Pattern) != "" {
		if strs.TrimSpace(props.Trigger+props.TriggerType) != "" || props.TriggerNegate || props.TriggerIgnoreCase {
			return sandbox.Deps.Std.Errorf("--pattern declares the args itself and excludes --trigger, --trigger-type, --trigger-negate and --trigger-ignore-case")
		}
		compiled, err := utils.CompileCommandPattern(sandbox, props.Pattern)
		if err != nil {
			return err
		}
		conf.Args = compiled.Args
		conf.Segments, conf.HasSegments = compiled.Segments, compiled.HasSegments
	} else {
		end := 0
		trigger := commandconf.Trigger{Exists: true, Type: "prefix", Values: []string{}}
		if props.Middleware && strs.TrimSpace(props.Trigger) == "" {
			end = commandconf.LastSegment
			if strs.TrimSpace(props.TriggerType) != "" || props.TriggerNegate || props.TriggerIgnoreCase {
				return sandbox.Deps.Std.Errorf("a --middleware with no --trigger runs on every command line: --trigger-type, --trigger-negate and --trigger-ignore-case need a --trigger")
			}
		} else {
			trigger_value, trigger_type := identifier, props.TriggerType
			if strs.TrimSpace(props.Trigger) != "" {
				trigger_value = props.Trigger
			}
			if props.Middleware && trigger_type == "" {
				trigger_type = "prefix"
			}
			built, err := utils.NewTrigger(sandbox, utils.TriggerProps{
				Type:       trigger_type,
				Value:      trigger_value,
				Negate:     props.TriggerNegate,
				IgnoreCase: props.TriggerIgnoreCase,
			})
			if err != nil {
				return err
			}
			trigger = built
		}
		if !trigger.Exists {
			return sandbox.Deps.Std.Errorf("a command needs a trigger on its first arg")
		}
		switch trigger.Type {
		case "equal":
			end = len(strs.Fields(trigger.Value)) - 1
		case "one-of":
		default:
			end = commandconf.LastSegment
		}
		conf.Args = []commandconf.Arg{{Id: "Command", Start: 0, End: end, Type: commandconf.DefaultArgType, Trigger: trigger}}
	}

	priority, err := addCommandPriority(sandbox, io, props)
	if err != nil {
		return err
	}
	conf.Priority = priority
	conf.Strict = !props.Middleware

	conf.Category = strs.TrimSpace(props.Category)
	if conf.Category == "" {
		conf.Category = defaultCategory
		if props.Middleware {
			conf.Category = defaultMiddlewareCategory
		}
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("add-command creating sandbox/internal/commands/%s \n", pkg)

	vars := map[string]interface{}{
		"Identifier":  identifier,
		"Package":     pkg,
		"Module":      module_conf.Module,
		"ProjectName": projectName(sandbox, io),
		"Pattern":     conf.Pattern(),
	}

	dir := utils.CommandDir(sandbox, name)
	if err := io.WriteFile(dir+"/"+utils.CommandConfFile, []byte(conf.Render())); err != nil {
		return err
	}

	template := CommandHandlerTemplate
	if props.Middleware {
		template = MiddlewareHandlerTemplate
	}
	handler, err := sandbox.Deps.Embeddeps.RenderTemplate(template, vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/"+utils.CommandHandlerFile, handler)
}

// addCommandPriority is the rung a new command lands on: --priority, one rung
// from the command --before or --after names, or the default of its kind.
func addCommandPriority(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddCommandProps) (int, error) {
	relative, has_relative, err := utils.CommandRelativePriority(sandbox, io, props.Before, props.After)
	if err != nil {
		return 0, err
	}
	if has_relative && props.HasPriority {
		return 0, sandbox.Deps.Std.Errorf("--priority excludes --before and --after")
	}
	switch {
	case has_relative:
		return relative, nil
	case props.HasPriority:
		if props.Priority < 0 {
			return 0, sandbox.Deps.Std.Errorf("--priority %d is negative: the chain runs from zero upwards", props.Priority)
		}
		return props.Priority, nil
	case props.Middleware:
		return commandconf.DefaultMiddlewarePriority, nil
	}
	return commandconf.DefaultPriority, nil
}

// CommandHandlerTemplate is the stub InternalPureHandler.go of a command, and
// MiddlewareHandlerTemplate the one of a middleware, which answers nothing.
const CommandHandlerTemplate = "templates/command_internal_pure_handler.go"
const MiddlewareHandlerTemplate = "templates/command_middleware_handler.go"

// projectName title-cases the target project's configured name for use in the
// scaffold's help text, falling back to the CLI's own ProjectName constant.
func projectName(sandbox *api.Sandbox, io *smartio.SmartIO) string {
	content, err := io.ReadFile(sandbox.Config.ProjectName + "Config/project.yaml")
	if err != nil {
		return sandbox.Config.ProjectName
	}
	conf, err := projectconf.New(sandbox, string(content))
	if err != nil || conf.Name == "" {
		return sandbox.Config.ProjectName
	}
	return sandbox.Deps.Stringsdeps.ToUpper(conf.Name[:1]) + conf.Name[1:]
}
