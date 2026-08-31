package actions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func ListAvailableExtensions(deps deps.Deps, props api.ListAvailableExtensionsProps) ([]string, error) {
	deps.Printf("list available extensions started with path %s \n", props.Path)
	return []string{}, nil
}
