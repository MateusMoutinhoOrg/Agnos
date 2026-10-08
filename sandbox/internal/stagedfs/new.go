package stagedfs

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/pathsconf"
)

func joinPath(sandbox *api.Sandbox, base string, name string) string {
	if sandbox.Deps.StringsDeps.HasSuffix(base, "/") || sandbox.Deps.StringsDeps.HasSuffix(base, "\\") {
		return base + name
	}
	return base + "/" + name
}

// normalizeRoot collapses the spellings of "the current directory" ("", ".",
// "./") to "" so rootedPath adds no prefix, and strips a trailing slash from
// every other value so joins are uniform.
func normalizeRoot(sandbox *api.Sandbox, path string) string {
	if path == "" || path == "." || path == "./" {
		return ""
	}
	for sandbox.Deps.StringsDeps.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return path
}

func New(sandbox *api.Sandbox, path string, projectName string) *StagedFS {
	io := &StagedFS{
		Root:         normalizeRoot(sandbox, path),
		sandbox:      sandbox,
		Transactions: make(map[string][]byte),
	}

	configDir := joinPath(sandbox, path, "Config")
	if projectName != "" {
		configDir = joinPath(sandbox, path, sandbox.Deps.StringsDeps.ToUpper(projectName[:1])+projectName[1:]+"Config")
	}

	replacersPath := joinPath(sandbox, configDir, "paths.yaml")
	if sandbox.Deps.IoDeps.Exists(replacersPath) && sandbox.Deps.IoDeps.IsFile(replacersPath) {
		content, err := sandbox.Deps.IoDeps.ReadFile(replacersPath)
		if err == nil {
			conf, err := pathsconf.New(sandbox, string(content))
			if err == nil {
				io.Replacers = conf
			} else {
				io.Replacers = pathsconf.NewEmpty(sandbox)
			}
		} else {
			io.Replacers = pathsconf.NewEmpty(sandbox)
		}
	} else {
		io.Replacers = pathsconf.NewEmpty(sandbox)
	}

	BindMethods(sandbox, io)
	return io
}
