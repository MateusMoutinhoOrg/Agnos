package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectBinds lists sandbox/binds and returns one "<Title>Bind" entry per
// .go file, for the {{range .Binds}} loop in sandbox/new.go.
func CollectBinds(deps *deps.Deps, io *smartio.SmartIO) []string {

	files := io.ListFiles("sandbox/binds")

	var binds []string
	for _, file := range files {
		parts := deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]

		if !deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}

		baseName := deps.Stringsdeps.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := deps.Stringsdeps.ToUpper(baseName[:1]) + baseName[1:]
		binds = append(binds, title+"Bind")
	}

	return binds
}
