package userconfig

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

type ModuleConf struct {
	Module    string
	GoVersion string
	Requires  []string

	Persist func() error
}

func NewModuleConf(sandbox *api.SandBox, path string) (*ModuleConf, error) {
	var module string
	var goversion string
	var requires []string

	if sandbox.Deps.IoLib.IsFile(path) {
		content_bytes, err := sandbox.Deps.IoLib.ReadFile(path)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(content_bytes), "\n")
		inRequireBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "//") {
				continue
			}

			if strings.HasPrefix(trimmed, "module ") {
				module = strings.TrimSpace(strings.TrimPrefix(trimmed, "module "))
			} else if strings.HasPrefix(trimmed, "go ") {
				goversion = strings.TrimSpace(strings.TrimPrefix(trimmed, "go "))
			} else if strings.HasPrefix(trimmed, "require (") {
				inRequireBlock = true
			} else if trimmed == ")" && inRequireBlock {
				inRequireBlock = false
			} else if inRequireBlock {
				requires = append(requires, trimmed)
			} else if strings.HasPrefix(trimmed, "require ") {
				req := strings.TrimSpace(strings.TrimPrefix(trimmed, "require "))
				requires = append(requires, req)
			}
		}
	}

	conf := &ModuleConf{
		Module:    module,
		GoVersion: goversion,
		Requires:  requires,
	}

	conf.Persist = func() error {
		var builder strings.Builder
		if conf.Module != "" {
			builder.WriteString("module " + conf.Module + "\n\n")
		}
		if conf.GoVersion != "" {
			builder.WriteString("go " + conf.GoVersion + "\n\n")
		}

		if len(conf.Requires) > 0 {
			builder.WriteString("require (\n")
			for _, req := range conf.Requires {
				builder.WriteString("\t" + req + "\n")
			}
			builder.WriteString(")\n")
		}

		return sandbox.Deps.IoLib.WriteFile(path, []byte(builder.String()))
	}

	return conf, nil
}