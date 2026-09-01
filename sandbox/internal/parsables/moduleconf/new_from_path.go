package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewFromPath(deps *deps.Deps, path string) (*ModuleConf, error) {
	bytes, err := deps.IoLib.ReadFile(path)
	if err != nil {
		return nil, deps.Errorf("failed to read module file at %s: %v", path, err)
	}

	return New(deps, string(bytes))
}
