package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectConstructors lists sandbox/api and returns one title-cased entry per
// .go file other than sandbox.go, for the {{range .Constructors}} loop in
// sandbox/api/sandbox.go. command.go and route.go are skipped along with it:
// what they declare is the shape of one command and of one route, not a field
// of the sandbox — the Commands and Routes slices built from them are written
// by the template itself.
func CollectConstructors(sandbox *api.Sandbox, io *smartio.SmartIO) []string {

	files := io.ListFiles("sandbox/api")

	var constructors []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]

		if name == "sandbox.go" || name == "command.go" || name == "route.go" || !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}

		baseName := sandbox.Deps.Stringsdeps.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := sandbox.Deps.Stringsdeps.ToUpper(baseName[:1]) + baseName[1:]
		constructors = append(constructors, title)
	}

	return constructors
}
