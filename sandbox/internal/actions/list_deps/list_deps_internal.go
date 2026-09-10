package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListDepsInternal returns the name of every dep of the embedded catalog, one
// entry per assets/deplist/ sub-directory, in listing order.
func ListDepsInternal(deps *deps.Deps, io *smartio.SmartIO, path string) ([]string, error) {
	deps.Std.Log("list-deps started with path %s \n", path)

	return utils.CatalogDeps(deps)
}
