package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddCommandInternal writes the two hand-written files of a new command
// package: a command.yaml whose one arg answers to the command's name on
// segment 0, and the InternalPureHandler.go stub. It refuses to overwrite an
// existing command (via io.WriteFile).
func AddCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string, help string, category string) error {
	if sandbox.Deps.Stringsdeps.TrimSpace(help) == "" {
		return sandbox.Deps.Std.Errorf("add-command requires --help")
	}
	if sandbox.Deps.Stringsdeps.TrimSpace(category) == "" {
		return sandbox.Deps.Std.Errorf("add-command requires --category")
	}

	if err := utils.ValidateCommandName(sandbox, name); err != nil {
		return err
	}
	utils.NoteNormalizedCommandName(sandbox, name)

	identifier := utils.CommandIdentifier(sandbox, name)
	pkg := utils.CommandPackage(sandbox, name)

	if utils.IsGeneratedCommand(sandbox, name) {
		return sandbox.Deps.Std.Errorf("the %s command is generated and cannot be declared", identifier)
	}

	sandbox.Deps.Std.Log("add-command creating sandbox/internal/commands/%s \n", pkg)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	conf := commandconf.NewEmpty(sandbox)
	conf.Args = []commandconf.Arg{{
		Id:      "Command",
		Start:   0,
		End:     0,
		Type:    commandconf.DefaultArgType,
		Trigger: commandconf.Trigger{Exists: true, Type: "equal", Value: identifier, Values: []string{}},
	}}
	conf.Category = sandbox.Deps.Stringsdeps.TrimSpace(category)
	conf.Help = sandbox.Deps.Stringsdeps.TrimSpace(help)

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

	handler, err := sandbox.Deps.Embeddeps.RenderTemplate(CommandHandlerTemplate, vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/"+utils.CommandHandlerFile, handler)
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
