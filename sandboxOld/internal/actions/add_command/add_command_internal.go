package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddCommandInternal writes the two hand-written files of a new command
// package. It refuses to overwrite an existing command (via io.WriteFile).
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

	if pkg == "help" {
		return sandbox.Deps.Std.Errorf("the help command is generated and cannot be declared")
	}

	sandbox.Deps.Std.Log("add-command creating sandbox/internal/commands/%s \n", pkg)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Identifier":  identifier,
		"Package":     pkg,
		"Module":      module_conf.Module,
		"ProjectName": projectName(sandbox, io),
		"Help":        sandbox.Deps.Stringsdeps.TrimSpace(help),
		"Category":    sandbox.Deps.Stringsdeps.TrimSpace(category),
	}

	dir := utils.CommandDir(sandbox, name)

	entries, err := sandbox.Deps.Embeddeps.RenderTemplate("templates/command_entries.yaml", vars)
	if err != nil {
		return err
	}
	if err := io.WriteFile(dir+"/entries.yaml", entries); err != nil {
		return err
	}

	handler, err := sandbox.Deps.Embeddeps.RenderTemplate("templates/command_handler.go", vars)
	if err != nil {
		return err
	}
	if err := io.WriteFile(dir+"/handler.go", handler); err != nil {
		return err
	}

	return nil
}

// projectName title-cases the target project's configured name for use in the
// scaffold's help text, falling back to the CLI's own ProjectName constant.
func projectName(sandbox *api.Sandbox, io *smartio.SmartIO) string {
	content, err := io.ReadFile(config.ProjectName + "Config/project.yaml")
	if err != nil {
		return config.ProjectName
	}
	conf, err := projectconf.New(sandbox, string(content))
	if err != nil || conf.Name == "" {
		return config.ProjectName
	}
	return sandbox.Deps.Stringsdeps.ToUpper(conf.Name[:1]) + conf.Name[1:]
}
