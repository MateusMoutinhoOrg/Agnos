package smartio

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/ignorableconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/pathreplacerconf"
)

func joinPath(base string, name string) string {
	if strings.HasSuffix(base, "/") || strings.HasSuffix(base, "\\") {
		return base + name
	}
	return base + "/" + name
}

// normalizeRoot collapses the spellings of "the current directory" ("", ".",
// "./") to "" so rootedPath adds no prefix, and strips a trailing slash from
// every other value so joins are uniform.
func normalizeRoot(path string) string {
	if path == "" || path == "." || path == "./" {
		return ""
	}
	for strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return path
}

func New(deps *deps.Deps, path string, projectName string) *SmartIO {
	io := &SmartIO{
		Root:         normalizeRoot(path),
		Transactions: make(map[string][]byte),
	}

	configDir := joinPath(path, projectName+"Config")

	ignorePath := joinPath(configDir, "ignore.yaml")
	if deps.Iodeps.Exist(ignorePath) && deps.Iodeps.IsFile(ignorePath) {
		content, err := deps.Iodeps.ReadFile(ignorePath)
		if err == nil {
			conf, err := ignorableconf.New(deps, string(content))
			if err == nil {
				io.Ignore = conf
			} else {
				io.Ignore = ignorableconf.NewEmpty(deps)
			}
		} else {
			io.Ignore = ignorableconf.NewEmpty(deps)
		}
	} else {
		io.Ignore = ignorableconf.NewEmpty(deps)
	}

	replacersPath := joinPath(configDir, "paths.yaml")
	if deps.Iodeps.Exist(replacersPath) && deps.Iodeps.IsFile(replacersPath) {
		content, err := deps.Iodeps.ReadFile(replacersPath)
		if err == nil {
			conf, err := pathreplacerconf.New(deps, string(content))
			if err == nil {
				io.Replacers = conf
			} else {
				io.Replacers = pathreplacerconf.NewEmpty(deps)
			}
		} else {
			io.Replacers = pathreplacerconf.NewEmpty(deps)
		}
	} else {
		io.Replacers = pathreplacerconf.NewEmpty(deps)
	}

	BindMethods(deps, io)
	return io
}
