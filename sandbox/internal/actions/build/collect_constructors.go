package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectConstructors lists sandbox/api and returns one title-cased entry per
// .go file other than sandbox.go, for the {{range .Constructors}} loop in
// sandbox/api/sandbox.go.
func CollectConstructors(deps *deps.Deps, io *smartio.SmartIO) []string {

	files := io.ListFiles("sandbox/api")

	var constructors []string
	for _, file := range files {
		parts := deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]

		if name == "sandbox.go" || !deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}

		baseName := deps.Stringsdeps.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := deps.Stringsdeps.ToUpper(baseName[:1]) + baseName[1:]
		constructors = append(constructors, title)
	}

	return constructors
}
