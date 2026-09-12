package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	serializables "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// RemoteModule is one resolved module of the Go module graph: the path it is
// imported under, the version that was resolved, and the directory holding its
// source — the module cache for a downloaded one, the replacement directory
// for one under a `replace`.
type RemoteModule struct {
	Path    string
	Version string
	Dir     string
}

// resolveModule finds a module's source on this machine, without writing
// anything. It asks the build list first (`go list -m`), which is what answers
// for a module already required — including one under a `replace`, the way a
// pair of repos is developed side by side — and falls back to downloading the
// named version.
func resolveModule(sandbox *api.Sandbox, path string, module string, version string) (*RemoteModule, error) {

	if resolved, err := readModule(sandbox, path, []string{"list", "-m", "-json", module}); err == nil {
		if version == "" || resolved.Version == version {
			return resolved, nil
		}
	}

	spec := module + "@latest"
	if version != "" {
		spec = module + "@" + version
	}

	return readModule(sandbox, path, []string{"mod", "download", "-json", spec})
}

// readModule runs one go command that prints a module as json and reads the
// three fields the installer needs out of it.
func readModule(sandbox *api.Sandbox, path string, args []string) (*RemoteModule, error) {
	result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     path,
		Program: "go",
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	if result.ExitCode != 0 {
		return nil, sandbox.Deps.Std.Errorf("go %s failed: %s", sandbox.Deps.Stringsdeps.Join(args, " "), result.Output)
	}

	parsed, err := sandbox.Deps.Serializables.ParseJson(result.Output)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("go %s printed something that is not json: %s", sandbox.Deps.Stringsdeps.Join(args, " "), result.Output)
	}

	module := &RemoteModule{
		Path:    jsonString(sandbox, parsed, "Path"),
		Version: jsonString(sandbox, parsed, "Version"),
		Dir:     jsonString(sandbox, parsed, "Dir"),
	}

	if module.Dir == "" {
		return nil, sandbox.Deps.Std.Errorf("go %s reported no source directory for the module", sandbox.Deps.Stringsdeps.Join(args, " "))
	}

	return module, nil
}

// jsonString reads one string field of a parsed json object, "" when it is
// absent or of another type.
func jsonString(sandbox *api.Sandbox, object *serializables.SerializibleObject, key string) string {
	item, _ := object.GetObjectItem(key)
	if item == nil || item.IsNull() {
		return ""
	}
	value, err := item.GetString()
	if err != nil {
		return ""
	}
	return value
}
