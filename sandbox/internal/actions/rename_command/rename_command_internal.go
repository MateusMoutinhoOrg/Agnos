package rename_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/commandconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatedCommandFiles are the files of a command package every build
// writes, so a rename leaves them behind for the follow-up build to write
// again.
var generatedCommandFiles = []string{"new.go", "entries.go"}

// RenameCommandInternal moves every hand-written file of
// sandbox/internal/commands/<command>/ to sandbox/internal/commands/<name>/,
// rewriting the package clause of each Go file, and removes the old
// directory. A verb the command answered to because it was its name — the
// equal trigger `add-command` wrote — becomes the new name; any other trigger
// is left as it is.
//
// The commands the build writes itself are refused.
func RenameCommandInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RenameCommandProps) error {
	command := utils.ResolveCommandName(sandbox, io, props.Command)
	if err := utils.ValidateCommandName(sandbox, command); err != nil {
		return err
	}
	if err := utils.ValidateCommandName(sandbox, props.Name); err != nil {
		return err
	}
	if utils.IsGeneratedCommand(sandbox, command) || utils.IsGeneratedCommand(sandbox, props.Name) {
		return sandbox.Deps.Std.Errorf("help, version and help-flag are generated and cannot be renamed")
	}

	old_pkg := utils.CommandPackage(sandbox, command)
	new_pkg := utils.CommandPackage(sandbox, props.Name)
	old_dir := utils.CommandDir(sandbox, command)
	new_dir := utils.CommandDir(sandbox, props.Name)

	if !io.IsDir(old_dir) {
		return sandbox.Deps.Std.Errorf("command %q not found in %s", utils.CommandIdentifier(sandbox, props.Command), old_dir)
	}
	if old_pkg == new_pkg {
		return sandbox.Deps.Std.Errorf("command %q is already named %q", utils.CommandIdentifier(sandbox, command), utils.CommandIdentifier(sandbox, props.Name))
	}
	if io.IsDir(new_dir) {
		return sandbox.Deps.Std.Errorf("command %q already exists in %s", utils.CommandIdentifier(sandbox, props.Name), new_dir)
	}

	conf, err := utils.LoadCommandConf(sandbox, io, command)
	if err != nil {
		return err
	}
	renameVerb(sandbox, conf, utils.CommandIdentifier(sandbox, command), utils.CommandIdentifier(sandbox, props.Name))

	sandbox.Deps.Std.Log("rename-command moving %s to %s \n", old_dir, new_dir)

	for _, file := range io.ListFilesRecursively(old_dir) {
		relative := sandbox.Deps.Stringsdeps.TrimPrefix(file, old_dir+"/")
		if isGenerated(relative) {
			continue
		}
		content, err := io.ReadFile(file)
		if err != nil {
			return err
		}
		if relative == utils.CommandConfFile {
			content = []byte(conf.Render())
		}
		if sandbox.Deps.Stringsdeps.HasSuffix(relative, ".go") {
			content = []byte(renamePackage(sandbox, string(content), old_pkg, new_pkg))
		}
		if err := io.WriteFile(new_dir+"/"+relative, content); err != nil {
			return err
		}
	}

	for _, file := range io.ListAllRecursively(old_dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(old_dir)
	return nil
}

// renameVerb points the equal trigger on segment 0 that spelled the old name
// at the new one.
func renameVerb(sandbox *api.Sandbox, conf *commandconf.CommandConf, old_verb string, new_verb string) {
	for index, arg := range conf.Args {
		if arg.Start == 0 && arg.Trigger.Exists && arg.Trigger.Type == "equal" && arg.Trigger.Value == old_verb {
			conf.Args[index].Trigger.Value = new_verb
			return
		}
	}
}

// isGenerated reports a file of the package the follow-up build writes again.
func isGenerated(relative string) bool {
	for _, generated := range generatedCommandFiles {
		if relative == generated {
			return true
		}
	}
	return false
}

// renamePackage rewrites the package clause of one Go file, and nothing else.
func renamePackage(sandbox *api.Sandbox, content string, old_pkg string, new_pkg string) string {
	lines := sandbox.Deps.Stringsdeps.Split(content, "\n")
	for i, line := range lines {
		if line == "package "+old_pkg {
			lines[i] = "package " + new_pkg
			break
		}
	}
	return sandbox.Deps.Stringsdeps.Join(lines, "\n")
}
