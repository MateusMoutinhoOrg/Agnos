package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*ModuleConf, error) {
	var module string
	var goversion string
	var requires []string
	var directives []string

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	lines := sandbox.Deps.Stringsdeps.Split(content, "\n")
	inRequireBlock := false
	inOtherBlock := false
	for _, line := range lines {
		trimmed := sandbox.Deps.Stringsdeps.TrimSpace(line)

		if inOtherBlock {
			directives = append(directives, line)
			if trimmed == ")" {
				inOtherBlock = false
			}
			continue
		}

		if trimmed == "" || sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "//") {
			continue
		}

		if sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "module ") {
			module = sandbox.Deps.Stringsdeps.TrimSpace(sandbox.Deps.Stringsdeps.TrimPrefix(trimmed, "module "))
		} else if sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "go ") {
			goversion = sandbox.Deps.Stringsdeps.TrimSpace(sandbox.Deps.Stringsdeps.TrimPrefix(trimmed, "go "))
		} else if sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "require (") {
			inRequireBlock = true
		} else if trimmed == ")" && inRequireBlock {
			inRequireBlock = false
		} else if inRequireBlock {
			requires = append(requires, trimmed)
		} else if sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "require ") {
			req := sandbox.Deps.Stringsdeps.TrimSpace(sandbox.Deps.Stringsdeps.TrimPrefix(trimmed, "require "))
			requires = append(requires, req)
		} else {
			// Every other directive is kept as it was written: this parser
			// models require alone, and what it does not model it must not
			// throw away.
			directives = append(directives, line)
			if sandbox.Deps.Stringsdeps.HasSuffix(trimmed, "(") {
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

	BindMethods(sandbox, conf)
	return conf, nil
}
