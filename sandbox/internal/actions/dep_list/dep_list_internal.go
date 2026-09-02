package dep_list

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// DepListInternal returns the name of every dep available under
// assets/deplist/, one entry per immediate sub-directory, in listing order.
func DepListInternal(deps *deps.Deps, io *smartio.SmartIO, path string) ([]string, error) {
	deps.Std.Printf("dep-list started with path %s \n", path)

	files, err := deps.Embeddeps.ListFilesRecursively("deplist")
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	deplist := []string{}
	for _, file := range files {
		name := strings.Split(file, "/")[0]
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		deplist = append(deplist, name)
	}

	return deplist, nil
}
