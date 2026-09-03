package add_command

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/projectconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

// AddCommandInternal writes the two hand-written files of a new command
// package. It refuses to overwrite an existing command (via io.WriteFile).
func AddCommandInternal(deps *deps.Deps, io *smartio.SmartIO, name string, help string, category string) error {
	if strings.TrimSpace(help) == "" {
		return deps.Std.Errorf("add-command requires --help")
	}
	if strings.TrimSpace(category) == "" {
		return deps.Std.Errorf("add-command requires --category")
	}

	identifier := utils.CommandIdentifier(name)
	pkg := utils.CommandPackage(name)
	if pkg == "" {
		return deps.Std.Errorf("invalid command name %q", name)
	}

	deps.Std.Printf("add-command creating sandbox/internal/commands/%s \n", pkg)

	gomod, err := io.ReadFile("go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Identifier":  identifier,
		"Package":     pkg,
		"Module":      module_conf.Module,
		"ProjectName": projectName(deps, io),
		"Help":        strings.TrimSpace(help),
		"Category":    strings.TrimSpace(category),
	}

	dir := utils.CommandDir(name)

	entries, err := deps.Embeddeps.RenderTemplate("templates/command_entries.yaml", vars)
	if err != nil {
		return err
	}
	if err := io.WriteFile(dir+"/entries.yaml", entries); err != nil {
		return err
	}

	handler, err := deps.Embeddeps.RenderTemplate("templates/command_handler.go", vars)
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
func projectName(deps *deps.Deps, io *smartio.SmartIO) string {
	content, err := io.ReadFile(config.ProjectName + "Config/project.yaml")
	if err != nil {
		return config.ProjectName
	}
	conf, err := projectconf.New(deps, string(content))
	if err != nil || conf.Name == "" {
		return config.ProjectName
	}
	return strings.ToUpper(conf.Name[:1]) + conf.Name[1:]
}
