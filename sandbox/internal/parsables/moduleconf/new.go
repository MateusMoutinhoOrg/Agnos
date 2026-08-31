package moduleconf

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*ModuleConf, error) {
	var module string
	var goversion string
	var requires []string

	if content == "" {
		return nil, deps.Errorf("content cannot be empty, use NewEmpty instead")
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

	BindMethods(deps, conf)
	return conf, nil
}
