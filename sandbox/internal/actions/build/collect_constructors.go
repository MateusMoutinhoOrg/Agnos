package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectConstructors lists sandbox/api and returns one title-cased entry per
// .go file other than sandbox.go, for the {{range .Constructors}} loop in
// sandbox/api/sandbox.go.
func CollectConstructors(sandbox *api.Sandbox, io *smartio.SmartIO) []string {

	files := io.ListFiles("sandbox/api")

	var constructors []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]

		if name == "sandbox.go" || !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") {
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
