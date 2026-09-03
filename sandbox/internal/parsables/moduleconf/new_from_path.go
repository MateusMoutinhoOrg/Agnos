package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func NewFromPath(deps *deps.Deps, path string) (*ModuleConf, error) {
	bytes, err := deps.Iodeps.ReadFile(path)
	if err != nil {
		return nil, deps.Std.Errorf("failed to read module file at %s: %v", path, err)
	}

	return New(deps, string(bytes))
}
