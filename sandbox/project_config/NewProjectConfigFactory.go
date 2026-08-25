package project_config

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewProjectConfigFactory(sandbox *api.SandBox) func(path string) api.ProjectConfig {
	return func(path string) api.ProjectConfig {

		projectConfig := api.ProjectConfig{}

		return projectConfig
	}
}
