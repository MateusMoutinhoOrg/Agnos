package project_config

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func SetProjectConfigMethods(sandbox *api.SandBox) {
	sandbox.NewProjectConfig = NewProjectConfigFactory(sandbox)

}
