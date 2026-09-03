package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// collectLibDirs lists the immediate sub-directories of dir and returns one
// {"Name": <dir>, "Title": <Dir>} entry per sub-directory, in listing order.
// Both the sandbox/deps and adapters/libs trees have this fixed shape: one
// sub-directory per sub-contract.
func collectLibDirs(io *smartio.SmartIO, dir string) []map[string]string {

	dirs := io.ListDirs(dir)

	var libs []map[string]string
	for _, d := range dirs {
		parts := strings.Split(d, "/")
		name := parts[len(parts)-1]

		if len(name) == 0 {
			continue
		}

		title := strings.ToUpper(name[:1]) + name[1:]
		libs = append(libs, map[string]string{
			"Name":  name,
			"Title": title,
		})
	}

	return libs
}

// CollectDepsLibs returns one entry per sandbox/deps sub-contract directory,
// for the {{range .DepsLibs}} loop in sandbox/deps/deps.go.
func CollectDepsLibs(io *smartio.SmartIO) []map[string]string {
	return collectLibDirs(io, "sandbox/deps")
}
