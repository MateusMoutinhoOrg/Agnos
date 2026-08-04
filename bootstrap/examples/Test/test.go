package main

import (
	bootstrapadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/adapters/standard"
	bootstraplib "github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox"
)

func main() {
	// 1. Build deps via the bootstrap adapter. The adapter is the only place
	//    allowed to reach outside the sandbox, so it is what initializes the
	//    embedded Agnos cache library (with the embedded library's own adapter).
	deps := bootstrapadapter.New("bootstrap-kvcache.json")

	// 2. Inject deps into the pure bootstrap library.
	l := bootstraplib.New(deps)

	// 3. The lib uses the embedded library through its Deps, never importing it.
	l.TestFunc()
}
