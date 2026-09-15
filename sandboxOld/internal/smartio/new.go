package smartio

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/ignorableconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/pathreplacerconf"
)

func joinPath(sandbox *api.Sandbox, base string, name string) string {
	if sandbox.Deps.Stringsdeps.HasSuffix(base, "/") || sandbox.Deps.Stringsdeps.HasSuffix(base, "\\") {
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
	for sandbox.Deps.Stringsdeps.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return path
}

func New(sandbox *api.Sandbox, path string, projectName string) *SmartIO {
	io := &SmartIO{
		Root:         normalizeRoot(sandbox, path),
		sandbox:      sandbox,
		Transactions: make(map[string][]byte),
	}

	configDir := joinPath(sandbox, path, projectName+"Config")

	ignorePath := joinPath(sandbox, configDir, "ignore.yaml")
	if sandbox.Deps.Iodeps.Exist(ignorePath) && sandbox.Deps.Iodeps.IsFile(ignorePath) {
		content, err := sandbox.Deps.Iodeps.ReadFile(ignorePath)
		if err == nil {
			conf, err := ignorableconf.New(sandbox, string(content))
			if err == nil {
				io.Ignore = conf
			} else {
				io.Ignore = ignorableconf.NewEmpty(sandbox)
			}
		} else {
			io.Ignore = ignorableconf.NewEmpty(sandbox)
		}
	} else {
		io.Ignore = ignorableconf.NewEmpty(sandbox)
	}

	replacersPath := joinPath(sandbox, configDir, "paths.yaml")
	if sandbox.Deps.Iodeps.Exist(replacersPath) && sandbox.Deps.Iodeps.IsFile(replacersPath) {
		content, err := sandbox.Deps.Iodeps.ReadFile(replacersPath)
		if err == nil {
			conf, err := pathreplacerconf.New(sandbox, string(content))
			if err == nil {
				io.Replacers = conf
			} else {
				io.Replacers = pathreplacerconf.NewEmpty(sandbox)
			}
		} else {
			io.Replacers = pathreplacerconf.NewEmpty(sandbox)
		}
	} else {
		io.Replacers = pathreplacerconf.NewEmpty(sandbox)
	}

	BindMethods(sandbox, io)
	return io
}
