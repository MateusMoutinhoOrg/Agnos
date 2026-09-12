package moduleconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewFromPath(sandbox *api.Sandbox, path string) (*ModuleConf, error) {
	bytes, err := sandbox.Deps.Iodeps.ReadFile(path)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("failed to read module file at %s: %v", path, err)
	}

	return New(sandbox, string(bytes))
}
