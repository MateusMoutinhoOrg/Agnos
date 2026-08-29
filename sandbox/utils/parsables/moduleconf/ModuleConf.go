package moduleconf

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

type ModuleConf struct {
	Module    string
	GoVersion string
	Requires  []string

	Render func() string
}

func addModuleConfMethods(sandbox *sandbox.SandBox, conf *ModuleConf) {
	conf.Render = func() string {
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

		return builder.String()
	}
}

func NewModuleConf(sandbox *sandbox.SandBox, content string) (*ModuleConf, error) {
	var module string
	var goversion string
	var requires []string

	if content == "" {
		return nil, sandbox.Deps.Errorf("content cannot be empty, use NewModuleConfEmpty instead")
	}

	lines := strings.Split(content, "\n")
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

	conf := &ModuleConf{
		Module:    module,
		GoVersion: goversion,
		Requires:  requires,
	}

	addModuleConfMethods(sandbox, conf)
	return conf, nil
}

func NewModuleConfEmpty(sandbox *sandbox.SandBox) *ModuleConf {
	conf := &ModuleConf{
		Requires: []string{},
	}
	addModuleConfMethods(sandbox, conf)
	return conf
}
