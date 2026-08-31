package actions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func RemoveExtension(deps *deps.Deps, props api.UninstallProps) error {
	deps.Printf("remove extension started with item %s \n", props.Item)
	return nil
}
