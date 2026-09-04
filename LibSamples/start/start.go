package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func main() {

	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
	module := "Test"
	startProps := api.StartProps{
		Path:        "blackbox",
		ProjectName: "Test",
		Module:      &module,
		Force:       false,
	}
	err := lib.Actions.Start(startProps)
	if err != nil {
		panic(err)
	}

}
