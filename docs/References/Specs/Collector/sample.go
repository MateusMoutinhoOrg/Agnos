package build

// This file is an illustrative sample, not part of the build.

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// DepsLib is one sub-contract directory: Name is the directory as written,
// Title the Deps field it becomes. Both are precomputed so the template only
// formats.
type DepsLib struct {
	Name  string
	Title string
}

// CollectDepsLibs lists sandbox/deps/<x>/ and returns one DepsLib per
// directory, in listing order, for {{range .DepsLibs}} in
// assets/deps/sandbox/deps/deps.go.
func CollectDepsLibs(io *smartio.SmartIO) []DepsLib {
	var libs []DepsLib
	for _, dir := range io.ListDirs("sandbox/deps") {
		parts := strings.Split(dir, "/")
		name := parts[len(parts)-1]
		libs = append(libs, DepsLib{
			Name:  name,
			Title: strings.ToUpper(name[:1]) + name[1:],
		})
	}
	return libs
}
