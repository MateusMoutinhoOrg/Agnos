package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
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
func resolveModule(deps *deps.Deps, path string, module string, version string) (*RemoteModule, error) {

	if resolved, err := readModule(deps, path, []string{"list", "-m", "-json", module}); err == nil {
		if version == "" || resolved.Version == version {
			return resolved, nil
		}
	}

	spec := module + "@latest"
	if version != "" {
		spec = module + "@" + version
	}

	return readModule(deps, path, []string{"mod", "download", "-json", spec})
}

// readModule runs one go command that prints a module as json and reads the
// three fields the installer needs out of it.
func readModule(deps *deps.Deps, path string, args []string) (*RemoteModule, error) {
	result, err := deps.Rundeps.Run(rundeps.RunProps{
		Dir:     path,
		Program: "go",
		Args:    args,
	})
	if err != nil {
		return nil, err
	}
	if result.ExitCode != 0 {
		return nil, deps.Std.Errorf("go %s failed: %s", deps.Stringsdeps.Join(args, " "), result.Output)
	}

	parsed, err := deps.Serializables.ParseJson(result.Output)
	if err != nil {
		return nil, deps.Std.Errorf("go %s printed something that is not json: %s", deps.Stringsdeps.Join(args, " "), result.Output)
	}

	module := &RemoteModule{
		Path:    jsonString(deps, parsed, "Path"),
		Version: jsonString(deps, parsed, "Version"),
		Dir:     jsonString(deps, parsed, "Dir"),
	}

	if module.Dir == "" {
		return nil, deps.Std.Errorf("go %s reported no source directory for the module", deps.Stringsdeps.Join(args, " "))
	}

	return module, nil
}

// jsonString reads one string field of a parsed json object, "" when it is
// absent or of another type.
func jsonString(deps *deps.Deps, object *serializables.SerializibleObject, key string) string {
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
