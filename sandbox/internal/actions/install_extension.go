package actions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func InstallExtension(deps *deps.Deps, props api.InstallProps) error {
	deps.Printf("install extension started with item %s \n", props.Item)
	return nil
}
