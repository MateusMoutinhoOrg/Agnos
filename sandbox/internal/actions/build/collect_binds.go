package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectBinds lists sandbox/binds and returns one "<Title>Bind" entry per
// .go file, for the {{range .Binds}} loop in sandbox/new.go.
func CollectBinds(io *smartio.SmartIO) []string {

	files := io.ListFiles("sandbox/binds")

	var binds []string
	for _, file := range files {
		parts := strings.Split(file, "/")
		name := parts[len(parts)-1]

		if !strings.HasSuffix(name, ".go") {
			continue
		}

		baseName := strings.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := strings.ToUpper(baseName[:1]) + baseName[1:]
		binds = append(binds, title+"Bind")
	}

	return binds
}
