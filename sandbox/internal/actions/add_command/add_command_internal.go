package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/triggerconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// defaultCategory is the heading a command lands under when it names none,
// and defaultMiddlewareCategory the one of a --middleware.
const defaultCategory = "Commands"
const defaultMiddlewareCategory = "Middleware"

// AddCommandInternal writes the two hand-written files of a new command
// package: a command.yaml and the handler.go stub. Its args are
// one arg answering to the command's name on segment 0 — every command line,
// for a --middleware — or another --trigger, or what a --pattern compiles to.
// It lands in the folder of its category under sandbox/internal/commands
// (Core -> commands/core/<name>) or the one props.Dir names instead, and
// refuses a name another command already carries, in whatever folder.
func AddCommandInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddCommandProps) error {
	strs := sandbox.Deps.StringsDeps
	help := strs.TrimSpace(props.Summary)
	if help == "" {
		return sandbox.Deps.StdDeps.Errorf("add-command requires --summary")
	}

	name := props.Name
	if err := utils.ValidateCommandName(sandbox, name); err != nil {
		return err
	}
	utils.NoteNormalizedCommandName(sandbox, name)

	identifier := utils.CommandName(sandbox, name)
	pkg := utils.CommandPackage(sandbox, name)

	if utils.IsGeneratedCommand(sandbox, name) {
		return sandbox.Deps.StdDeps.Errorf("the %s command is generated and cannot be declared", identifier)
	}

	category := strs.TrimSpace(props.Category)
	if category == "" {
		category = defaultCategory
		if props.Middleware {
			category = defaultMiddlewareCategory
		}
	}

	// The folder is the category's, so the path says which heading the
	// command is listed under; --dir overrides it.
	folder := props.Dir
	if strs.TrimSpace(folder) == "" {
		folder = category
	}
	group, err := utils.UnitGroup(sandbox, folder)
	if err != nil {
		return err
	}
	if existing, found := utils.FindUnitDir(sandbox, io, utils.CommandsDir, utils.CommandConfFile, pkg); found {
		return sandbox.Deps.StdDeps.Errorf("command %q already exists in %s: a command name is unique across every folder", identifier, existing)
	}
	dir := utils.UnitDirIn(utils.CommandsDir, group, pkg)

	conf := commandconf.NewEmpty(sandbox)
	conf.Summary = help

	if strs.TrimSpace(props.Pattern) != "" {
		if strs.TrimSpace(props.Trigger+props.TriggerType) != "" || props.TriggerNegate || props.TriggerIgnoreCase {
			return sandbox.Deps.StdDeps.Errorf("--pattern declares the args itself and excludes --trigger, --trigger-type, --trigger-negate and --trigger-ignore-case")
		}
		compiled, err := utils.CompileCommandPattern(sandbox, props.Pattern)
		if err != nil {
			return err
		}
		conf.Args = compiled.Args
		conf.Segments, conf.HasSegments = compiled.Segments, compiled.HasSegments
	} else {
		end := 0
		trigger := triggerconf.Trigger{Set: true, Type: "prefix", Values: []string{}}
		if props.Middleware && strs.TrimSpace(props.Trigger) == "" {
			end = commandconf.LastSegment
			if strs.TrimSpace(props.TriggerType) != "" || props.TriggerNegate || props.TriggerIgnoreCase {
				return sandbox.Deps.StdDeps.Errorf("a --middleware with no --trigger runs on every command line: --trigger-type, --trigger-negate and --trigger-ignore-case need a --trigger")
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
		if !trigger.Set {
			return sandbox.Deps.StdDeps.Errorf("a command needs a trigger on its first arg")
		}
		switch trigger.Type {
		case "equal":
			end = len(strs.Fields(trigger.Value)) - 1
		case "one-of":
		default:
			end = commandconf.LastSegment
		}
		conf.Args = []commandconf.Arg{{Id: "Verb", Start: 0, End: end, Type: commandconf.DefaultArgType, Trigger: trigger}}
	}

	priority, err := addCommandPriority(sandbox, io, props)
	if err != nil {
		return err
	}
	conf.Priority = priority
	conf.Strict = !props.Middleware

	conf.Category = category

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("add-command creating %s \n", dir)

	vars := map[string]interface{}{
		"Identifier":  identifier,
		"Package":     pkg,
		"Module":      module_conf.Module,
		"ProjectName": projectName(sandbox, io),
		"Pattern":     conf.Pattern(),
	}

	if err := io.CreateFile(dir+"/"+utils.CommandConfFile, []byte(conf.Render())); err != nil {
		return err
	}

	template := CommandHandlerTemplate
	if props.Middleware {
		template = MiddlewareHandlerTemplate
	}
	handler, err := sandbox.Deps.EmbedDeps.RenderTemplate(template, vars)
	if err != nil {
		return err
	}
	return io.CreateFile(dir+"/"+utils.CommandHandlerFile, handler)
}

// addCommandPriority is the rung a new command lands on: --priority, one rung
// from the command --before or --after names, or the default of its kind.
func addCommandPriority(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddCommandProps) (int, error) {
	relative, has_relative, err := utils.CommandRelativePriority(sandbox, io, props.Before, props.After)
	if err != nil {
		return 0, err
	}
	if has_relative && props.HasPriority {
		return 0, sandbox.Deps.StdDeps.Errorf("--priority excludes --before and --after")
	}
	switch {
	case has_relative:
		return relative, nil
	case props.HasPriority:
		if props.Priority < 0 {
			return 0, sandbox.Deps.StdDeps.Errorf("--priority %d is negative: the chain runs from zero upwards", props.Priority)
		}
		return props.Priority, nil
	case props.Middleware:
		return commandconf.DefaultMiddlewarePriority, nil
	}
	return commandconf.DefaultPriority, nil
}

// CommandHandlerTemplate is the stub handler.go of a command, and
// MiddlewareHandlerTemplate the one of a middleware, which answers nothing.
const CommandHandlerTemplate = "templates/command_handler.go"
const MiddlewareHandlerTemplate = "templates/command_middleware_handler.go"

// projectName is the target project's declared name for use in the
// scaffold's help text, falling back to the CLI's own ProjectName.
func projectName(sandbox *api.Sandbox, io *stagedfs.StagedFS) string {
	content, err := io.ReadFile(utils.ProjectConfPath(sandbox))
	if err != nil {
		return sandbox.Config.ProjectName
	}
	conf, err := projectconf.New(sandbox, string(content))
	if err != nil || conf.ProjectName == "" {
		return sandbox.Config.ProjectName
	}
	return conf.ProjectName
}
