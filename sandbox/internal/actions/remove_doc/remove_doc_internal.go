package remove_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveDocInternal deletes every file under the doc's directory plus the
// directory itself, sub-docs and assets included. The generated docs/Index/
// is not a doc and is refused: build rewrites it whole.
func RemoveDocInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateDocName(sandbox, name); err != nil {
		return err
	}

	dir := utils.DocDir(sandbox, name)
	if !io.IsDir(dir) {
		return sandbox.Deps.Std.Errorf("doc %s not found", dir)
	}

	sandbox.Deps.Std.Log("remove-doc removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
