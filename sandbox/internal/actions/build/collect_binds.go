package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectBinds lists sandbox/binds and returns one "<Title>Bind" entry per
// .go file, for the {{range .Binds}} loop in sandbox/new.go.
func CollectBinds(sandbox *api.Sandbox, io *smartio.SmartIO) []string {

	files := io.ListFiles("sandbox/binds")

	var binds []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]

		if !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}

		baseName := sandbox.Deps.Stringsdeps.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := sandbox.Deps.Stringsdeps.ToUpper(baseName[:1]) + baseName[1:]
		binds = append(binds, title+"Bind")
	}

	return binds
}
