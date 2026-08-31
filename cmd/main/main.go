package main

import (
	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnosapi "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"

	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {

	deps := agnosadapter.New(".agnos")

	l := agnoslib.New(&deps)
	l.Build(agnosapi.BuildProps{
		Path: ".",
	})

}
