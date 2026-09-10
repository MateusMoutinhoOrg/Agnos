package dep_list

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DepListInternal returns the name of every dep of the embedded catalog, one
// entry per assets/deplist/ sub-directory, in listing order.
func DepListInternal(deps *deps.Deps, io *smartio.SmartIO, path string) ([]string, error) {
	deps.Std.Log("dep-list started with path %s \n", path)

	return utils.CatalogDeps(deps)
}
