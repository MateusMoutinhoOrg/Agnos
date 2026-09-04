package remove_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveDocInternal deletes every file under the doc's directory plus the
// directory itself, sub-docs and assets included. The generated docs/Index/
// is not a doc and is refused: build rewrites it whole.
func RemoveDocInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateDocName(deps, name); err != nil {
		return err
	}

	dir := utils.DocDir(name)
	if !io.IsDir(dir) {
		return deps.Std.Errorf("doc %s not found", dir)
	}

	deps.Std.Log("remove-doc removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
