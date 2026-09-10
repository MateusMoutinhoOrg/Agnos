package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*ModuleConf, error) {
	var module string
	var goversion string
	var requires []string
	var directives []string

	if content == "" {
		return nil, deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	lines := deps.Stringsdeps.Split(content, "\n")
	inRequireBlock := false
	inOtherBlock := false
	for _, line := range lines {
		trimmed := deps.Stringsdeps.TrimSpace(line)

		if inOtherBlock {
			directives = append(directives, line)
			if trimmed == ")" {
				inOtherBlock = false
			}
			continue
		}

		if trimmed == "" || deps.Stringsdeps.HasPrefix(trimmed, "//") {
			continue
		}

		if deps.Stringsdeps.HasPrefix(trimmed, "module ") {
			module = deps.Stringsdeps.TrimSpace(deps.Stringsdeps.TrimPrefix(trimmed, "module "))
		} else if deps.Stringsdeps.HasPrefix(trimmed, "go ") {
			goversion = deps.Stringsdeps.TrimSpace(deps.Stringsdeps.TrimPrefix(trimmed, "go "))
		} else if deps.Stringsdeps.HasPrefix(trimmed, "require (") {
			inRequireBlock = true
		} else if trimmed == ")" && inRequireBlock {
			inRequireBlock = false
		} else if inRequireBlock {
			requires = append(requires, trimmed)
		} else if deps.Stringsdeps.HasPrefix(trimmed, "require ") {
			req := deps.Stringsdeps.TrimSpace(deps.Stringsdeps.TrimPrefix(trimmed, "require "))
			requires = append(requires, req)
		} else {
			// Every other directive is kept as it was written: this parser
			// models require alone, and what it does not model it must not
			// throw away.
			directives = append(directives, line)
			if deps.Stringsdeps.HasSuffix(trimmed, "(") {
				inOtherBlock = true
			}
		}
	}

	conf := &ModuleConf{
		Module:     module,
		GoVersion:  goversion,
		Requires:   requires,
		Directives: directives,
	}

	BindMethods(deps, conf)
	return conf, nil
}
