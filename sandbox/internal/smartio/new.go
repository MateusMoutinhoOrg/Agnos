package smartio

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/ignorableconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/pathreplacerconf"
)

func joinPath(base string, name string) string {
	if strings.HasSuffix(base, "/") || strings.HasSuffix(base, "\\") {
		return base + name
	}
	return base + "/" + name
}

func New(deps *deps.Deps, path string, projectName string) *SmartIO {
	io := &SmartIO{
		Transactions: make(map[string][]byte),
	}

	configDir := joinPath(path, projectName+"Config")

	ignorePath := joinPath(configDir, "ignore.yaml")
	if deps.IoLib.Exist(ignorePath) && deps.IoLib.IsFile(ignorePath) {
		content, err := deps.IoLib.ReadFile(ignorePath)
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
	if deps.IoLib.Exist(replacersPath) && deps.IoLib.IsFile(replacersPath) {
		content, err := deps.IoLib.ReadFile(replacersPath)
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
